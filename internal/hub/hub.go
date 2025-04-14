// hub отвечает за контроль над пользователями и вебсокет-клиентами
package hub

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
	"github.com/anrew1002/Tournament-ChemLoto/internal/usecase"
	"github.com/gorilla/websocket"
)

type ChannelRepository interface {
	Get(channel string) ([]string, bool)
	Add(channel string, connection string)
	Remove(channel string, connection string)
	GetChannelFunc(channel string) (func(chan common.Message), bool)
	SetChannelFunc(channel string, fun func(chan common.Message))
}

type RoomRepository interface {
	AddRoom(name string, engine models.Engine) (*entities.Room, error)
	GetRooms() ([]*entities.Room, error)
	GetRoom(name string) (*entities.Room, error)
	SubscribeToRoom(name string, user *entities.User) error
}

type Hub struct {
	log *slog.Logger
	db  *sql.DB

	// Rooms    *roomsState
	Rooms2   RoomRepository
	Users2   *repository.UserRepository
	upgrader websocket.Upgrader

	Connections *connectionsState
	Channels    ChannelRepository
	Channels2   *repository.GroupsRepository

	// eventHandlers     map[string]HandlerFunc
	WebsocketHandlers *WebsocketHandlers
	eventChan         chan internalEventWrap
}

func NewHub(log *slog.Logger, uc *usecase.Usecases, upgrader websocket.Upgrader, db *sql.DB) *Hub {
	roomRepo := repository.NewRoomRepo(db)
	groupRepo := repository.NewGroupsRepo(db)
	userRepo := repository.NewUserRepo(db)
	wh := NewWebsocketHandlers(uc, log.With("origin: websocketHandlers"))
	return &Hub{
		upgrader:          upgrader,
		log:               log,
		Users2:            userRepo,
		Rooms2:            roomRepo,
		Connections:       &connectionsState{state: make(map[string]*SockConnection)},
		Channels:          repository.NewChannelState(),
		Channels2:         groupRepo,
		WebsocketHandlers: wh,
		// eventHandlers:     make(map[string]HandlerFunc),
		eventChan: make(chan internalEventWrap, 10),
	}
}

func (h *Hub) Run() {
	go func() {
	eventLoop:
		for event := range h.eventChan {
			h.log.Debug("finding handler...")
			// handler, ok := h.eventHandlers[event.msgType.String()]
			// if !ok {
			// 	h.log.Error("error find event handler ", "message type", event.msgType.String())
			// 	continue eventLoop
			// }
			// h.log.Debug("start handle event in Hub.Run", "event", event.msgType.String())
			// go handler.HandleEvent(h, event)

			err := h.WebsocketHandlers.Handle(event)
			if err != nil {
				h.log.Error("error find event handler ", "message type", event.msgType.String())
				continue eventLoop
			}

			h.log.Debug("exit handler")
		}
	}()
}
func (h *Hub) SendEventToHub(e internalEventWrap) {
	h.log.Debug("send event to Hub")
	h.eventChan <- e
}

// TODO:
// func (h *Hub) SaveGamesStats() map[string]*bytes.Buffer {
// 	h.Rooms.mutex.Lock()
// 	defer h.Rooms.mutex.Unlock()

// 	results := make(map[string]*bytes.Buffer)
// 	for k, v := range h.Rooms.state {
// 		// Инициализация нового bytes.Buffer
// 		buffer := new(bytes.Buffer)

// 		// Создание нового csv.Writer
// 		writer := csv.NewWriter(buffer)

// 		// Запись данных в CSV
// 		err := writer.WriteAll(v.Engine.GetResults())
// 		if err != nil {
// 			h.log.Error("Error writing to buffer:", slog.String("room", k), sl.Err(err))
// 			continue
// 		}

// 		// Сохранение буфера в результаты
// 		results[v.Name] = buffer
// 	}

//		return results
//	}

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
