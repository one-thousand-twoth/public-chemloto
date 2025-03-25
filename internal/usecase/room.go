package usecase

import (
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
)

type CreateRoomRequest struct {
	Name         string `json:"name" validate:"required,min=1,safeinput"`
	Type         string `json:"type"`
	EngineConfig map[string]any
}

type Response struct {
	Rooms any      `json:"rooms"`
	Error []string `json:"error"`
}

func CreateRoom(repo *repository.RoomRepository, req CreateRoomRequest, log *slog.Logger) (*entities.Room, error) {

	const op = "server.handlers.CreateRoom"

	eng, err := engines.NewEngine(req.Type, req.Name, log, req.EngineConfig)
	if err != nil {
		return nil, enerr.E(op, err)
	}

	room, err := repo.AddRoom(req.Name, eng)
	if err != nil {
		return nil, enerr.E(op, err)
	}

	return room, nil
}

func GetRooms(repo *repository.RoomRepository) ([]*entities.Room, error) {
	const op = "server.handlers.CreateRoom"

	rooms, err := repo.GetRooms()
	if err != nil {
		return nil, enerr.E(op, err)
	}

	return rooms, nil
}
