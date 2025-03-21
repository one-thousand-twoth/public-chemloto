package server

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/hub"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/golang-cz/devslog"
	"github.com/gorilla/websocket"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
)

type Server struct {
	mux  *chi.Mux
	db   *sql.DB
	code string
	hub  *hub.Hub
	log  *slog.Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run(port string) {
	QuitChannel := make(chan struct{})
	signalQuit := make(chan os.Signal, 1)
	admin_code := flag.String("code", "test_code", "password for accessing administrator")

	flag.Parse()
	s.code = *admin_code
	// log.Println("Starting server on port", port)
	// log.Fatal(http.ListenAndServe(":"+port, s))
	srv := http.Server{Addr: "0.0.0.0:" + port, Handler: s.mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("Failed listening", sl.Err(err))
			QuitChannel <- struct{}{}
		}
	}()
	ip := GetOutboundIP().String()
	s.log.Info("Сервер запущен:", slog.String("IP address", "http://"+ip+":"+port))
	s.log.Info("Cекретный код:",
		slog.String("Password", *admin_code))
	signal.Notify(signalQuit, syscall.SIGINT, syscall.SIGTERM)

	// Дальше смерть
	select {
	case <-QuitChannel:
	case <-signalQuit:
	}
	s.log.Info("Shutting down")
	doneChan := make(chan struct{})
	go func() {
		s.log.Info("Сохранение результатов игр...")
		results := s.hub.SaveGamesStats()
		for roomName, result := range results {
			filePath := filepath.Join("chemloto results", roomName+".csv")
			file, err := os.Create(filePath)
			if err != nil {
				s.log.Error("Failed to create result file", sl.Err(err))
				return
			}
			defer file.Close()

			_, err = result.WriteTo(file)
			if err != nil {
				s.log.Error("Failed to write in result file", sl.Err(err))
				return
			}
			s.log.Info("Saved results for room", "room", roomName)
		}
		doneChan <- struct{}{}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.log.Error("Error during server shutdown",
			sl.Err(err),
		)
	}

	<-doneChan

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

	log := slog.New(devslog.NewHandler(colorable.NewColorableStdout(), opts))
	mux := chi.NewRouter()

	db := sqlite.MustInitDB()

	Hub := hub.NewHub(log, upgrader, db)
	Hub.SetupHandlers()
	Hub.Run()
	// NOTE: for development
	Hub.FillRooms()

	server := &Server{
		hub: Hub,
		log: log,
		mux: mux,
		db:  db}
	server.configureRoutes()

	return server
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// log.Fatal("GetOutboundIP", "Fatal error", err)
		return net.ParseIP("0.0.0.0")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
