package server

import (
	"net/http"
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
		s.log.Info("team registred", "name", name)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(token))

	}
}
