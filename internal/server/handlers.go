package server

import (
	"log/slog"
	"net/http"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/anrew1002/Tournament-ChemLoto/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	valid "github.com/go-playground/validator/v10"
)

func (s *Server) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encode(w, r, http.StatusOK, map[string]any{
			"Channels":    s.hub.Channels,
			"Connections": s.hub.Connections,
			"Rooms":       s.hub.Rooms,
			"Users":       s.hub.Users,
		})
	}
}

func (s *Server) GetRooms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encode(w, r, http.StatusOK, s.hub.Rooms)
	}
}
func (s *Server) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encode(w, r, http.StatusOK, s.hub.Users)
	}
}
func (s *Server) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clnt, ok := s.hub.Users.GetByToken(token)
		if !ok {
			encode(w, r, http.StatusNotFound, clnt)
			return
		}
		encode(w, r, http.StatusOK, clnt)
	}
}
func (s *Server) CreateRoom() http.HandlerFunc {
	type Request struct {
		Name       string         `json:"name" validate:"required,min=3,safeinput"`
		MaxPlayers int            `json:"maxPlayers" validate:"required"`
		Elements   map[string]int `json:"elementCounts" validate:"required"`
		Time       int            `validate:"excluded_if=isAuto false,gte=0"`
		IsAuto     bool           `json:"isAuto"`
	}
	type Response struct {
		Rooms any      `json:"rooms"`
		Error []string `json:"error"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.CreateRoom"
		log := s.log.With(slog.String("op", op))
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body", sl.Err(err))
			encode(w, r, http.StatusBadRequest, Response{Error: []string{"Неправильный формат запроса"}})
			return
		}
		log.Debug("request body decoded", slog.Any("request", req))
		validate := validator.Ins
		if err := validate.Struct(req); err != nil {
			validateErr := err.(valid.ValidationErrors)
			encode(w, r, http.StatusBadRequest, Response{Error: validator.ValidationError(validateErr)})
			return
		}
		if !req.IsAuto {
			req.Time = 0
		}
		room := hub.NewRoom(req.Name, req.MaxPlayers, req.Elements, req.Time, req.IsAuto,
			polymers.New(
				s.log.With(slog.String("room", req.Name)),
				polymers.PolymersEngineConfig{
					Elements: req.Elements,
					TimerInt: req.Time,
					Unicast: func(userID string, msg common.Message) {
						s.log.Debug("Unicast message")
						usr, ok := s.hub.Users.Get(userID)
						if !ok {
							s.log.Error("failed to get user while Unicast message from engine")
						}
						connID := usr.GetConnection()
						conn, ok := s.hub.Connections.Get(connID)
						conn.MessageChan <- msg
					},
					Broadcast: func(msg common.Message) {
						s.log.Debug("Broadcast message")
						s.hub.SendMessageOverChannel(req.Name, msg)
					},
				},
			))

		// TODO: добавить обработку уже имеющейся комнаты
		if err := s.hub.Rooms.Add(room); err != nil {
			log.Error("failed to add room", sl.Err(err))
			encode(w, r, http.StatusConflict, Response{Error: []string{"Сервер не смог создать комнату"}})
			return
		}
		s.log.Info("Room created", "name", req.Name, "time", req.Time)
		// s.hub.SendMessageOverChannel("default", models.Message{Type: websocket.TextMessage, Body: []byte(req.Name)})
		encode(w, r, http.StatusOK, Response{Rooms: s.hub.Rooms, Error: []string{}})
	}
}

func (s *Server) Login(AdminCode string) http.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,min=3,safeinput"`
		Code string `json:"code,omitempty"`
	}
	type Response struct {
		Token string   `json:"token"`
		Error []string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.Login"
		log := s.log.With(slog.String("op", op))
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body")
			encode(w, r, http.StatusBadRequest, Response{Error: []string{"Неправильный формат запроса"}})
			return
		}
		log.Debug("request body decoded", slog.Any("request", req))
		validate := validator.Ins
		if err := validate.Struct(req); err != nil {
			validateErr := err.(valid.ValidationErrors)
			encode(w, r, http.StatusBadRequest, Response{Error: validator.ValidationError(validateErr)})
			return
		}
		var role hub.Role
		role = hub.Player_Role
		if req.Code == AdminCode {
			role = hub.Admin_Role
		}

		token, err := GenerateRandomStringURLSafe(32)
		if err != nil {
			encode(w, r, http.StatusInternalServerError, Response{Error: []string{"Ошибка сервера"}})
			return
		}
		channels := []string{"default"}
		user := hub.NewUser(req.Name, token, "", role, channels)
		if err := s.hub.Users.Add(user); err != nil {
			encode(w, r, http.StatusConflict, Response{Error: []string{"Пользователь с таким именем уже существует"}})
			return
		}
		s.log.Info("user registred", "name", req.Name)
		encode(w, r, http.StatusCreated, Response{Token: token})
	}
}
