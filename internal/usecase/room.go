package usecase

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
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

func (uc *Usecases) CreateRoom(req CreateRoomRequest, log *slog.Logger) (*entities.Room, error) {

	const op = "server.handlers.CreateRoom"

	eng, err := engines.NewEngine(
		req.Type,
		req.Name,
		log,
		req.EngineConfig,
		func(userID string, msg common.Message) {
			go func() {
				user, err := uc.userRepo.GetUserByID(entities.ID(32))
				if err != nil {
					log.Error("Error getting user while unicast")
					return
				}
				user.MessageChan <- msg
			}()
		},
		func(msg common.Message) {
			go func() {
				user, err := uc.userRepo.GetUserByID(entities.ID(32))
				if err != nil {
					log.Error("Error getting user while unicast")
					return
				}
				user.MessageChan <- msg
			}()
		},
	)

	if err != nil {
		return nil, enerr.E(op, err)
	}

	room, err := uc.roomRepo.AddRoom(req.Name, eng)
	if err != nil {
		return nil, enerr.E(op, err)
	}

	return room, nil
}

func (uc *Usecases) GetRooms(ctx context.Context) ([]*entities.Room, error) {
	rows, err := uc.roomRepo.GetRooms()

	if err != nil {
		return nil, err
	}

	rooms := make([]*entities.Room, 0, len(rows))
	for _, v := range rows {
		rooms = append(rooms, &entities.Room{
			Name:   v.Name,
			Engine: v.Engine,
		})
	}

	return rooms, err

}

func (uc *Usecases) SubscribeToRoom(ctx context.Context, roomName string, userID entities.ID) error {
	const op enerr.Op = "usecase.room/SubscribeToRoom"

	tx, err := uc.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	queries := uc.queries.WithTx(tx)

	rowUser, err := queries.GetUserByID(ctx, int64(userID))
	if err != nil {
		return enerr.E(op, err)
	}

	user := entities.ToUserModel(rowUser)

	err = user.SubscribeToRoom(roomName)
	if err != nil {
		return enerr.E(op, err)
	}

	err = queries.UpdateUserRoom(ctx, database.UpdateUserRoomParams{
		Room: sql.NullString{
			String: user.Room,
			Valid:  true,
		},
		ID: int64(user.ID),
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (uc *Usecases) StartGame(
	userID entities.ID,
) error {
	const op enerr.Op = "usecase.room/StartGame"

	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return enerr.E(op, err)
	}

	if !user.HasPermision() {
		return enerr.E("No permission to start game")
	}

	room, err := uc.roomRepo.GetRoom(user.Room)
	if err != nil {
		return err
	}

	room.Engine.Start()

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
