package server

import (
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
	s.mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/rooms", s.GetRoom())
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
