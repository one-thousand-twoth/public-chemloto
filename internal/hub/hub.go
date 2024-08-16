// hub отвечает за контроль над пользователями и вебсокет-клиентами
package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
)

type Hub struct {
	log *slog.Logger

	Rooms    *roomsState
	Users    *usersState
	upgrader websocket.Upgrader

	Connections *connectionsState
	Channels    *channelsState

	eventHandlers map[string]HandlerFunc
	eventChan     chan internalEventWrap
}

func NewHub(log *slog.Logger, upgrader websocket.Upgrader) *Hub {
	return &Hub{
		upgrader:    upgrader,
		log:         log,
		Rooms:       &roomsState{state: make(map[string]*room)},
		Users:       &usersState{state: make(map[string]*User)},
		Connections: &connectionsState{state: make(map[string]*SockConnection)},
		Channels: &channelsState{
			state: make(map[string]mapset.Set[string]),
			initFunctions: map[string]func(chan common.Message){
				"default": func(c chan common.Message) {},
			},
		},
		eventHandlers: make(map[string]HandlerFunc),
		eventChan:     make(chan internalEventWrap, 10),
	}
}

func (h *Hub) Run() {
	go func() {
	eventLoop:
		for event := range h.eventChan {
			handler, ok := h.eventHandlers[event.msgType.String()]
			if !ok {
				h.log.Error("error find event handler ", "message type", event.msgType.String())
				continue eventLoop
			}
			h.log.Debug("start handle event in Hub.Run", "event", event.msgType.String())
			handler.HandleEvent(h, event)
		}
	}()
}
func (h *Hub) SendEventToHub(e internalEventWrap) {
	h.eventChan <- e
}

func (h *Hub) SendMessageOverChannel(channel string, message common.Message) error {
	op := "SendMessageOverChannel"
	log := h.log.With("op", op)
	connections, ok := h.Channels.Get(channel)
	if !ok {
		return fmt.Errorf("no channel with %s name", channel)
	}
	log.Debug(fmt.Sprintf("WS Message: %+v to %s", message, channel))
	for _, connectionName := range connections {
		conn, ok := h.Connections.Get(connectionName)
		if ok {
			conn.MessageChan <- message
		}
	}
	return nil
}

func GetMessageType(msg []byte) (payload map[string]interface{}, msgType common.MessageType, err error) {

	err = json.Unmarshal(msg, &payload)
	if err != nil {
		err = errors.New("fail to unmarshal struct")
		return
	}
	payloadType, ok := payload["Type"]
	if !ok {
		err = errors.New("поле Type не указано")
		return
	}
	msgType, ok = common.MapEnumStringToMessageType[payloadType.(string)]
	if !ok {
		err = errors.New("неизвестный тип сообщения")
		return
	}
	return payload, msgType, err
}

func (h *Hub) CheckToken(token string) (*User, error) {
	if token == "" {
		return nil, fmt.Errorf("bad token")
	}
	clnt, ok := h.Users.GetByToken(token)
	if !ok {
		return nil, fmt.Errorf("bad token")
	}
	return clnt, nil
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	op := "handlers.handleWS"
	log := h.log.With("op", op)
	reqToken := r.URL.Query().Get("token")
	if reqToken == "" {
		http.Error(w, "Incorrect token", http.StatusBadRequest)
		return
	}
	user, err := h.CheckToken(reqToken)
	if err != nil {
		log.Error(op, sl.Err(err))
		http.Error(w, "Bad Token", http.StatusUnauthorized)
		return
	}
	user.mutex.Lock()
	defer user.mutex.Unlock()

	username := user.Name
	log = log.With("userWS", username)
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade connection", sl.Err(err))
		return
	}

	connection := NewConnection(conn, username)
	// (1) Добавляем соединение в общий список
	h.Connections.Add(connection)
	// h.Users.Add(client)]
	// TODO: действительно ли оно надо?
	// (2) Указываем новое соединение пользователю
	user.conn = connection.ID

	for _, channel := range user.channels {
		// (3) Добавляем к необходимым каналам новое соединение
		h.Channels.Add(channel, connection.ID)
		channel := channel
		// Вызываем Init функцию для канала, если есть
		go func() {
			initChan, ok := h.Channels.GetChannelFunc(channel)
			if !ok {
				return
			}
			log.Debug("init func start Hub")
			initChan(connection.MessageChan)
		}()
	}
	log.Debug("user", "user current channels:", fmt.Sprintf("%+v", user.channels))
	if x, ok := h.Channels.Get("default"); ok {
		log.Debug("hub current channels", "hub channels:", fmt.Sprintf("%+#v", x))
	}
	go func() {
	ReceiveLoop:
		for {
			messageType, msg, err := connection.Conn.ReadMessage()
			if err != nil {
				// Disconnect on error
				if connection.CloseChannel != nil {
					connection.CloseChannel <- struct{}{}
				}
				break
			}

			switch messageType {
			case websocket.CloseMessage:
				if connection.CloseChannel != nil {
					connection.CloseChannel <- struct{}{}
				}
				break ReceiveLoop
			case websocket.TextMessage:
				msg, msgType, err := GetMessageType(msg)
				if err != nil {
					log.Error("failed to get message type", "type", msgType)
				}
				log.Debug(fmt.Sprintf("got messageWS %+v", msg))
				h.SendEventToHub(NewEventWrap(username, user.Room, msg, msgType))

			// TODO: Реализовать ping на клиенте
			// Handling receiving ping/pong
			case websocket.PingMessage:
				fallthrough
			case websocket.PongMessage:
				// connection.lock.Lock()
				// connection.lastPing = time.Now()
				// connection.lock.Unlock()
			}
		}
	}()

	sendMutex := sync.Mutex{}
	go func() {
	SendLoop:
		for {
			select {
			case message := <-connection.MessageChan:
				envelope, err := json.Marshal(message)
				if err != nil {
					log.Error("Marshaling data", slog.Any("data", message))
				}
				log.Debug("envelope", "env", string(envelope))
				sendMutex.Lock()
				_ = conn.WriteMessage(websocket.TextMessage, envelope)
				sendMutex.Unlock()
			case <-connection.CloseChannel:
				log.Debug("Disconnecting client", "client", username)

				//Сигнализируем что соединение было закрыто
				connection.CloseChannel = nil

				sendMutex.Lock()

				// Send close message with 1000
				_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				// Sleep a tiny bit to allow message to be sent before closing connection
				time.Sleep(time.Millisecond)
				_ = conn.Close()
				// (1) Удаляем соединение из общего списка
				h.Connections.Remove(connection.ID)
				// (2) Удаляем соединение у пользователя
				user.SetConnection("")
				// (3) Удаляем из подписок каналов соединие
				for _, channel := range user.GetChannels() {
					h.Channels.Remove(channel, connection.ID)
				}

				break SendLoop
			}
		}
	}()
}
