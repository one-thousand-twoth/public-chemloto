package server

import (
	"net/http"

	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/go-chi/chi/v5"
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

		if err := s.hub.Users.Add(name, token); err != nil {
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
	roomID := chi.URLParam(r, "room")
	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := s.hub.Rooms.Get(roomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reqToken := r.URL.Query().Get("token")
	if reqToken == "" {
		http.Error(w, "Incorrect token", http.StatusBadRequest)
		return
	}
	client, err := s.CheckToken(reqToken)
	if err != nil {
		s.log.Error(op, sl.Err(err))
		http.Error(w, "Bad Token", http.StatusUnauthorized)
		return
	}
	s.log.Debug("client connected", "op", op, "name", client.Name)
	// w.Header().Set("Sec-Websocket-Protocol", reqToken)
	// s.upgrader.Subprotocols = []string{reqToken}
	// TODO:
	_, err = s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Failed to upgrade connection", sl.Err(err))
		return
	}

	// client.SetConnection(conn)
	// client.StartTicker(room.TickerDuration)
	// s.hub.Register(client, room)

	// s.hub.ListenClient(client, room)
}
