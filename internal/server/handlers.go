package server

import (
	"log/slog"
	"net/http"

	"github.com/anrew1002/Tournament-ChemLoto/internal/hub"
	"github.com/anrew1002/Tournament-ChemLoto/internal/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	valid "github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
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

func (s *Server) GetRoom() http.HandlerFunc {
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
		s.hub.SendMessageOverChannel("default", models.Message{Type: websocket.TextMessage, Body: []byte(name)})
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
		user := hub.NewUser(name, token, "", hub.Player_Role, channels)
		if err := s.hub.Users.Add(user); err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		s.log.Info("user registred", "name", name)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(token))

	}
}

func (s *Server) Login(AdminCode string) http.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,safeinput"`
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
