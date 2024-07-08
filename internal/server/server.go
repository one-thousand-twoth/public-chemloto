package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/hub"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Server struct {
	mux      *chi.Mux
	upgrader websocket.Upgrader
	hub      *hub.Hub
	log      *slog.Logger
	storage  *sqlite.Storage
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run(port string) {
	QuitChannel := make(chan struct{})
	signalQuit := make(chan os.Signal, 1)
	// log.Println("Starting server on port", port)
	// log.Fatal(http.ListenAndServe(":"+port, s))
	srv := http.Server{Addr: "0.0.0.0:" + port, Handler: s.mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("Failed listening", sl.Err(err))
			QuitChannel <- struct{}{}
		}
	}()

	s.log.Info("Listening...")
	signal.Notify(signalQuit, syscall.SIGINT, syscall.SIGTERM)

	// Дальше смерть
	select {
	case <-QuitChannel:
	case <-signalQuit:
	}
	s.log.Info("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.log.Error("Error during server shutdown",
			sl.Err(err),
		)
	}

	// Allow time to disconnect & clear from Redis
	time.Sleep(time.Second)

	s.log.Info("Stopped")
}

func NewServer() *Server {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{"token"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	mux := chi.NewRouter()
	storage := sqlite.NewStorage()

	hub := hub.NewHub(storage, log)
	// hub.Run()

	server := &Server{
		upgrader: upgrader,
		hub:      hub,
		log:      log,
		storage:  storage,
		mux:      mux}
	server.configureRoutes()

	return server
}
