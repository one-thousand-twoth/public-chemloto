package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/hub"
	"github.com/anrew1002/Tournament-ChemLoto/internal/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func (s *Server) GetRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encode(w, r, http.StatusOK, s.hub.Rooms)
	}
}

func (s *Server) CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var name string
		if name = r.Form.Get("name"); name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.hub.Rooms.Add(name); err != nil {
			encode(w, r, http.StatusConflict, s.hub.Rooms)
			return
		}
		encode(w, r, http.StatusOK, s.hub.Rooms)
	}
}

func (s *Server) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var name string
		if name = r.Form.Get("name"); name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token, err := GenerateRandomStringURLSafe(32)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		channels := []string{"default"}
		usr := hub.NewUser(name, token, nil, channels)
		if err := s.hub.Users.Add(usr); err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		s.log.Info("user registred", "name", name)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(token))

	}
}

func (s *Server) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clnt := s.hub.Users.Get(token)
		if clnt == nil {
			encode(w, r, http.StatusNotFound, clnt)
		}
		encode(w, r, http.StatusOK, clnt)

	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	op := "handlers.handleWS"

	reqToken := r.URL.Query().Get("token")
	if reqToken == "" {
		http.Error(w, "Incorrect token", http.StatusBadRequest)
		return
	}
	user, err := s.CheckToken(reqToken)
	if err != nil {
		s.log.Error(op, sl.Err(err))
		http.Error(w, "Bad Token", http.StatusUnauthorized)
		return
	}

	s.log.Debug("client connected", "op", op, "name", user.Name)
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Failed to upgrade connection", sl.Err(err))
		return
	}

	connection := hub.NewConnection(conn, user.Name)
	// (1) Добавляем соединение в общий список
	s.hub.Connections.Add(connection)
	// s.hub.Users.Add(client)]
	// TODO: действительно ли оно надо?
	// (2) Указываем новое соединение пользователю
	user.SetConnection(connection)

	for _, channel := range user.GetChannels() {
		// (3) Добавляем к необходимым каналам новое соединение
		s.hub.Channels.Add(channel, connection.ID)
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
				connection.MessageChan <- models.Message{Type: websocket.TextMessage, Body: msg}
				s.hub.Input(string(msg))
			// Handling receiving ping/pong
			case websocket.PingMessage:
				fallthrough
			case websocket.PongMessage:
				// connection.lock.Lock()
				// connection.lastPing = time.Now()
				// connection.lock.Unlock()
			}
			// //  определение типа сообщения
			// // fmt.Printf("InputUser: %s", string(payload))
			// payload, payloadType, err := getMessageType(msg)
			// if err != nil {
			// 	h.log.Debug("Fail get messageType", sl.Err(err))
			// 	h.RaiseWSError(err.Error(), client)
			// 	continue
			// }
			// h.log.Debug("get payload", "struct", payload)

			// switch payloadType {
			// case "Input":
			// 	h.InputGameLogic(payload, room, client)
			// }

		}
	}()

	sendMutex := sync.Mutex{}

SendLoop:
	for {
		select {
		case message := <-connection.MessageChan:
			sendMutex.Lock()
			_ = conn.WriteMessage(int(message.Type), message.Body)
			sendMutex.Unlock()
		case <-connection.CloseChannel:
			// logger.Info("Disconnecting user",
			// 	zap.String("requestId", requestid.Get(c)),
			// 	zap.String("id", connId),
			// )

			//Сигнализируем что соединение было закрыто
			connection.CloseChannel = nil

			sendMutex.Lock()

			// Send close message with 1000
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			// Sleep a tiny bit to allow message to be sent before closing connection
			time.Sleep(time.Millisecond)
			_ = conn.Close()
			// (1) Удаляем соединение из общего списка
			s.hub.Connections.Remove(connection.ID)
			// (2) Удаляем соединение у пользователя
			user.SetConnection(nil)
			// (3) Удаляем из подписок каналов соединие
			for _, channel := range user.GetChannels() {
				s.hub.Channels.Remove(channel, connection.ID)
			}

			break SendLoop
		}
	}
	// client.SetConnection(conn)
	// client.StartTicker(room.TickerDuration)
	// s.hub.Register(client, room)

	// s.hub.ListenClient(client, room)
}
