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
	// stripped, err := fs.Sub(web.DIST, "dist")
	// if err != nil {
	// 	s.log.Error("Failed get FS", sl.Err(err))
	// }

	// frontendFS := http.FileServer(http.FS(stripped))
	// s.mux.Handle("/*", frontendFS)
	s.mux.NotFound(NotFoundHandler)
	// s.mux.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	// 	if _, err := os.Stat(r.RequestURI); os.IsNotExist(err) {
	// 		http.StripPrefix(r.RequestURI, frontendFS).ServeHTTP(w, r)
	// 	} else {
	// 		frontendFS.ServeHTTP(w, r)
	// 	}
	// })

	// items := http.FileServer(http.FS(web.DIST))
	// s.mux.Handle("/*", items)

	s.mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/rooms", s.GetRooms())
		r.Post("/rooms", s.CreateRoom())
		// r.Post("/users", s.CreateUser())
		r.Get("/users", s.GetUsers())
		r.Post("/users", s.Login("test_code"))
		r.Get("/users/{token}", s.GetUser())
		r.Get("/ws", s.hub.HandleWS)
		r.Get("/debug", s.Status())

		// // Subrouters:
		// r.Route("/{articleID}", func(r chi.Router) {
		// 	r.Use(ArticleCtx)
		// 	r.Get("/", getArticle)       // GET /articles/123
		// 	r.Put("/", updateArticle)    // PUT /articles/123
		// 	r.Delete("/", deleteArticle) // DELETE /articles/123
		// })
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
