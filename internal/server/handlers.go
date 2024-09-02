package server

import (
	"log/slog"
	"net/http"

	"github.com/anrew1002/Tournament-ChemLoto/internal/appvalidation"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func (s *Server) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encode(w, r, http.StatusOK, map[string]any{
			"Users":       s.hub.Users,
			"Channels":    s.hub.Channels,
			"Rooms":       s.hub.Rooms,
			"Connections": s.hub.Connections,
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
func (s *Server) DeleteUser() http.HandlerFunc {
	type Response struct {
		Error []string `json:"error"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.PatchUser"
		log := s.log.With(slog.String("op", op))
		username := chi.URLParam(r, "username")
		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clnt, ok := s.hub.Users.Get(username)
		if !ok {
			encode(w, r, http.StatusNotFound, Response{Error: []string{"Возможно пользователь был удалён"}})
			return
		}
		if clnt.Room != "" {
			encode(w, r, http.StatusBadRequest, Response{Error: []string{"Нельзя изменить пользователя находящегося в игре"}})
			return
		}

		s.hub.Users.Remove(clnt.Name)

		log.Info("Deleted user", "user", username)

		encode(w, r, http.StatusOK, Response{})
	}
}
func (s *Server) PatchUser() http.HandlerFunc {
	type Request struct {
		Role common.Role
	}
	type Response struct {
		Error []string `json:"error"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.PatchUser"
		log := s.log.With(slog.String("op", op))
		username := chi.URLParam(r, "username")
		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body", sl.Err(err))
			encode(w, r, http.StatusBadRequest, Response{Error: []string{"Неправильный формат запроса"}})
			return
		}
		clnt, ok := s.hub.Users.Get(username)
		if !ok {
			encode(w, r, http.StatusNotFound, Response{Error: []string{"Возможно пользователь был удалён"}})
			return
		}
		if clnt.Room != "" {
			encode(w, r, http.StatusBadRequest, Response{Error: []string{"Нельзя изменить пользователя в игре"}})
			return
		}
		clnt.Role = req.Role
		encode(w, r, http.StatusOK, clnt)
	}
}
func (s *Server) CreateRoom() http.HandlerFunc {
	type Request struct {
		hub.Room
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
		// TODO: добавить обработку уже имеющейся комнаты
		if err := s.hub.AddNewRoom(req.Room); err != nil {
			log.Error("failed to add room", sl.Err(err))
			switch validateErr := err.(type) {
			case validator.ValidationErrors:
				encode(w, r, http.StatusBadRequest, Response{Error: appvalidation.ValidationError(validateErr)})
				return
			}
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
		Token string      `json:"token"`
		Role  common.Role `json:"role"`
		Error []string    `json:"error"`
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
		validate := appvalidation.Ins
		if err := validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			encode(w, r, http.StatusBadRequest, Response{Error: appvalidation.ValidationError(validateErr)})
			return
		}
		var role common.Role
		role = common.Player_Role
		if req.Code == AdminCode {
			role = common.Admin_Role
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
		s.log.Info("user registred", "name", req.Name, "role", role)
		encode(w, r, http.StatusCreated, Response{Token: token, Role: role})
	}
}
