package hub

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/gorilla/websocket"
)

func (h *Hub) HandleWS2(w http.ResponseWriter, r *http.Request) {
	op := "handlers.handleWS"

	log := h.log.With("op", op)
	reqToken := r.URL.Query().Get("token")

	if reqToken == "" {
		http.Error(w, "StatusUnauthorized", http.StatusBadRequest)
		return
	}

	user, err := h.usecases.UserRepo.GetUserByApikey(reqToken)
	if err != nil {
		http.Error(w, "StatusUnauthorized", http.StatusUnauthorized)
		return
	}
	if user.MessageChan == nil {
		panic("NIL CHANNEL!")
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade connection", sl.Err(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	connection := NewConnection(conn, user.Room)

	channels, err := h.usecases.GroupRepo.GetAllUserGroups(user.ID)
	if err != nil {
		channels = make([]*entities.Group, 0)
	}

	for _, channel := range channels {
		channel := channel
		// Вызываем Init функцию для канала, если есть
		go func() {
			log.Debug(fmt.Sprintf("Running initFunction for %s", channel))
			channel.Fn(user.MessageChan)
		}()
	}

	go func() {
		recieve(connection, log, user, h)
	}()

	go func() {
		send(connection, log, conn, user)
	}()
	user.MessageChan <- common.Message{
		Type:   common.ENGINE_ACTION,
		Ok:     true,
		Errors: []string{},
		Body:   map[string]any{"message": "hello! What are you doing there?"},
	}
	log.Debug("User connected to websocket", "user", user)
}

func send(connection *SockConnection, log *slog.Logger, conn *websocket.Conn, user *entities.User) {
	sendMutex := sync.Mutex{}
SendLoop:
	for {
		select {
		case message := <-user.MessageChan:
			envelope, err := json.Marshal(message)
			if err != nil {
				log.Error("Marshaling data", slog.Any("data", message))
			}
			log.Debug("envelope", "env", string(envelope))
			sendMutex.Lock()
			_ = conn.WriteMessage(websocket.TextMessage, envelope)
			sendMutex.Unlock()
		case <-connection.CloseChannel:
			log.Debug("Disconnecting client", "client", user.Name)

			//Сигнализируем что соединение было закрыто
			connection.CloseChannel = nil

			sendMutex.Lock()

			// Send close message with 1000
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			// Sleep a tiny bit to allow message to be sent before closing connection
			time.Sleep(time.Millisecond)
			_ = conn.Close()

			break SendLoop
		}
	}
}

func recieve(connection *SockConnection, log *slog.Logger, user *entities.User, h *Hub) {
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
			// user.mutex.Lock()
			h.SendEventToHub(
				NewEventWrap(
					user.ID,
					msg,
					msgType,
					user.MessageChan,
				))
			// user.mutex.Unlock()
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
}
