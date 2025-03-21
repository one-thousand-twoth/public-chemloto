package server

import (
	"embed"
	"errors"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/anrew1002/Tournament-ChemLoto/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) configureRoutes() {
	s.mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	s.mux.Mount("/debug", middleware.Profiler())
	s.mux.NotFound(NotFoundHandler)
	s.mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/rooms", s.GetRooms())
		r.Post("/rooms", s.CreateRoom())
		r.Get("/users", s.GetUsers())
		r.Post("/users/{username}", s.PatchUser())
		r.Delete("/users/{username}", s.DeleteUser())
		r.Post("/users", s.Login())
		r.Get("/users/{token}", s.GetUser())
		r.Get("/ws", s.hub.HandleWS)
		r.Get("/debug", s.Status())

	})
	s.mux.Route("/api/v2", func(r chi.Router) {
		r.Get("/rooms", s.GetRooms2())
		r.Post("/rooms", s.CreateRoom2())
		// r.Get("/users", s.GetUsers())
		// r.Post("/users/{username}", s.PatchUser())
		// r.Delete("/users/{username}", s.DeleteUser())
		r.Post("/users", s.Login2())
		r.Get("/users/{token}", s.GetUser2())
		r.Get("/ws", s.hub.HandleWS2)
		// r.Get("/debug", s.Status())

	})
}

var ErrDir = errors.New("path is dir")

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	err := tryRead(web.DIST, "dist", r.URL.Path, w)
	if err == nil {
		return
	}

	err = tryRead(web.DIST, "dist", "index.html", w)
	if err != nil {
		panic(err)
	}
}

func tryRead(fs embed.FS, prefix, requestedPath string, w http.ResponseWriter) error {
	f, err := fs.Open(path.Join(prefix, requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return ErrDir
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}
