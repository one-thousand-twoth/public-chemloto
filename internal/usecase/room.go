package usecase

import (
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
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

type RoomRepository interface {
	AddRoom(name string, engine models.Engine) (*entities.Room, error)
	GetRooms() ([]*entities.Room, error)
	GetRoom(name string) (*entities.Room, error)
	SubscribeToRoom(name string, user *entities.User) error
}

func CreateRoom(repo RoomRepository, req CreateRoomRequest, log *slog.Logger) (*entities.Room, error) {

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

func GetRooms(repo RoomRepository) ([]*entities.Room, error) {
	const op = "server.handlers.CreateRoom"

	rooms, err := repo.GetRooms()
	if err != nil {
		return nil, enerr.E(op, err)
	}

	return rooms, nil
}

func SubscribeToRoom(
	roomRepo RoomRepository,
	roomName string,
	user *entities.User,
) error {
	const op enerr.Op = "usecase.subscribtions/SubscribeToRoom"

	// if data.Target == "" || data.Name == "" {
	// 	return enerr.E(op, "empty field", enerr.Validation)
	// }

	err := roomRepo.SubscribeToRoom(roomName, user)
	if err != nil {
		return err
	}

	return nil

}

func StartGame(
	roomRepo RoomRepository,
	userRepo *repository.UserRepository,
	roomName string,
	userID entities.ID,
) error {
	const op enerr.Op = "usecase.subscribtions/StartGame"

	// if data.Target == "" || data.Name == "" {
	// 	return enerr.E(op, "empty field", enerr.Validation)
	// }

	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return enerr.E(op, err)
	}

	err = roomRepo.SubscribeToRoom(user.Room, user)
	if err != nil {
		return err
	}

	return nil

}

// func UnSubscribeToRoom(
// 	roomRepo RoomRepository,
// 	roomName string,
// 	user *entities.User,
// ) error {
// 	const op enerr.Op = "usecase.subscribtions/UnsubscribeToRoom"

// 	// if data.Target == "" || data.Name == "" {
// 	// 	return enerr.E(op, "empty field", enerr.Validation)
// 	// }

// 	err := roomRepo.Un(roomName, user)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
