package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/anrew1002/Tournament-ChemLoto/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (s *Server) Login() http.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,min=1,safeinput"`
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

		request := usecase.LoginRequest{
			Name: req.Name,
			Code: req.Code,
		}
		user, err := s.usecases.Login(s.log, request, s.code)
		if err != nil {
			encodeError(w, log, err)
			return
		}

		encode(w, r, http.StatusCreated, Response{Token: user.Token, Role: user.Role})
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
		err = s.usecases.PatchUserRole(context.TODO(), usecase.PatchRequest{
			Name: username,
			Role: req.Role,
		})
		if err != nil {
			log.Error("failed to PatchUserRole", sl.Err(err))
			encode(w, r, http.StatusBadRequest, Response{Error: []string{"Неправильный формат запроса"}})
		}

		encode(w, r, http.StatusOK, struct{}{})
	}
}

func (s *Server) CreateRoom() http.HandlerFunc {
	type Request struct {
		Name         string         `json:"name" validate:"required,min=1,safeinput"`
		Type         string         `json:"type"`
		EngineConfig map[string]any `json:"engineConfig"`
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
		createRoomParams := usecase.CreateRoomParams{
			Name:         req.Name,
			Type:         req.Type,
			EngineConfig: req.EngineConfig,
		}

		_, err = s.usecases.CreateRoom(createRoomParams, s.log)
		if err != nil {
			encodeError(w, log, err)
			return
		}
		encode(w, r, http.StatusOK, Response{Rooms: nil, Error: []string{}})
	}
}

func (s *Server) GetRooms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := s.usecases.GetRooms(context.TODO())
		if err != nil {
			encode(w, r, http.StatusInternalServerError, struct{}{})
			return
		}
		roomMap := make(map[string]*entities.Room, len(rooms))
		for _, v := range rooms {
			roomMap[v.Name] = v
		}
		encode(w, r, http.StatusOK, roomMap)
	}
}

func (s *Server) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clnt, err := s.usecases.UserRepo.GetUserByApikey(token)
		if err != nil {
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
		log.Error("User delete", "user", username)
		err := s.usecases.DeleteUser(context.TODO(), username)
		if err != nil {
			encodeError(w, log, err)
			return
		}

		log.Info("Deleted user", "user", username)
		encode(w, r, http.StatusOK, Response{})
	}
}

func (s *Server) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op enerr.Op = "server.handlers.GetUsers"
		log := s.log.With(enerr.OpAttr(op))
		users, err := s.usecases.GetUsers()
		if err != nil {
			encodeError(w, log, err)
		}
		encode(w, r, http.StatusOK, users)
	}
}
