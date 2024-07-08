package server

import "github.com/go-chi/chi/v5"

func (s *Server) configureRoutes() {
	s.mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/rooms", s.GetRoom())
		r.Post("/rooms", s.CreateRoom())

		// // Subrouters:
		// r.Route("/{articleID}", func(r chi.Router) {
		// 	r.Use(ArticleCtx)
		// 	r.Get("/", getArticle)       // GET /articles/123
		// 	r.Put("/", updateArticle)    // PUT /articles/123
		// 	r.Delete("/", deleteArticle) // DELETE /articles/123
		// })
	})
}
