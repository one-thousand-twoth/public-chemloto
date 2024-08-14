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
	"github.com/go-chi/chi/v5"
	"github.com/golang-cz/devslog"
	"github.com/gorilla/websocket"
)

type Server struct {
	mux *chi.Mux
	hub *hub.Hub
	log *slog.Logger
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
	slogOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	// new logger with options
	opts := &devslog.Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[15:04:05]",
		NewLineAfterLog:   true,
		DebugColor:        devslog.Magenta,
	}

	log := slog.New(devslog.NewHandler(os.Stdout, opts))
	mux := chi.NewRouter()

	hub := hub.NewHub(log, upgrader)
	hub.SetupHandlers()
	hub.Run()

	server := &Server{
		hub: hub,
		log: log,
		mux: mux}
	server.configureRoutes()

	return server
}
