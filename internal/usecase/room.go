package usecase

import (
	"context"
	"database/sql"
	"log/slog"
	"regexp"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/invopop/validation"
)

type CreateRoomParams struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	EngineConfig map[string]any `json:"engineConfig"`
}

func (arg *CreateRoomParams) Validate() error {
	return validation.ValidateStruct(&arg,
		validation.Field(&arg.Name,
			validation.Required,
			validation.Length(1, 25).Error("Должно быть меньше 26 символов"),
			validation.Match(regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9- ]+$`)).
				Error("Должно содержать только буквы, цифры и пробел"),
		),
		validation.Field(&arg.Name,
			validation.In("polymers", "amino").Error("Неизвестный тип движка"),
		),
		validation.Field(&arg.EngineConfig,
			validation.NotNil.Error("Конфиг не должен быть пустым"),
		),
	)
}

type Response struct {
	Rooms any      `json:"rooms"`
	Error []string `json:"error"`
}

func (uc *Usecases) CreateRoom(args CreateRoomParams, log *slog.Logger) (*entities.Room, error) {

	const op = "server.handlers.CreateRoom"

	log.Info("Creating room...", slog.Any("args", args))

	eng, err := engines.NewEngine(
		args.Type,
		args.Name,
		log,
		args.EngineConfig,
		func(username string, msg common.Message) {
			go func() {
				user, err := uc.UserRepo.GetUserByName(username)
				if err != nil {
					log.Error("Error getting user while unicast")
					return
				}
				user.MessageChan <- msg
			}()
		},
		func(msg common.Message) {
			go func() {
				users, err := uc.UserRepo.GetRoomSubscribers(args.Name)
				if err != nil {
					log.Error("Error getting user while broadcast", slog.Any("err", err.Error()))
					return
				}
				for _, user := range users {
					user.MessageChan <- msg
				}
			}()
		},
	)

	if err != nil {
		return nil, enerr.E(op, err)
	}

	room, err := uc.RoomRepo.AddRoom(args.Name, eng)
	if err != nil {
		return nil, enerr.E(op, err)
	}

	return room, nil
}

func (uc *Usecases) GetRooms(ctx context.Context) ([]*entities.Room, error) {
	rows, err := uc.RoomRepo.GetRooms()

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
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
	// TODO:
	engine, err := uc.RoomRepo.GetEngine(roomName)

	err = engine.AddParticipant(models.Participant{
		Name: user.Name,
		Role: user.Role,
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	// then NOTIFY
	users, err := uc.RoomRepo.GetRoomUsers(context.TODO(), roomName)
	if err != nil {
		// TODO:
		return nil
	}
	data := engine.PreHook()
	for _, v := range users {

		v.MessageChan <- common.Message{
			Type: common.ENGINE_INFO,
			Ok:   true,
			Body: data,
		}
	}

	return nil

}
func (uc *Usecases) UnsubscribeFromRoom(ctx context.Context, roomName string, userID entities.ID) error {
	const op enerr.Op = "usecase.room/UnsubscribeFromRoom"

	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := uc.queries.WithTx(tx)

	rowUser, err := queries.GetUserByID(ctx, int64(userID))
	if err != nil {
		return enerr.E(op, err)
	}

	user := entities.ToUserModel(rowUser)

	err = user.UnsubscribeFromRoom(roomName)
	if err != nil {
		return enerr.E(op, err)
	}

	err = queries.UpdateUserRoom(ctx, database.UpdateUserRoomParams{
		Room: sql.NullString{
			String: "",
			Valid:  false,
		},
		ID: int64(user.ID),
	})
	if err != nil {
		return err
	}
	// TODO:
	engine, _ := uc.RoomRepo.GetEngine(roomName)

	err = engine.RemoveParticipant(user.Name)
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

	user, err := uc.UserRepo.GetUserByID(userID)
	if err != nil {
		return enerr.E(op, err)
	}

	if !user.HasPermision() {
		return enerr.E("No permission to start game")
	}

	room, err := uc.RoomRepo.GetRoom(user.Room)
	if err != nil {
		return err
	}

	room.Engine.Start()

	return nil

}

func (uc *Usecases) GetRoomUsers(ctx context.Context, roomname string) ([]*entities.User, error) {
	const op enerr.Op = "usecase.room/GetRoomUsers"
	users, err := uc.RoomRepo.GetRoomUsers(ctx, roomname)
	if err != nil {
		return nil, enerr.E(op, err)
	}
	return users, nil
}

func (uc *Usecases) StopGame(ctx context.Context, roomName string, userID entities.ID) error {
	const op enerr.Op = "usecase.room/StopGame"

	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return enerr.E(op, err, enerr.Internal)
	}
	defer tx.Rollback()

	queries := uc.queries.WithTx(tx)

	rowUser, err := queries.GetUserByID(ctx, int64(userID))
	if err != nil {
		return enerr.E(op, err)
	}

	user := entities.ToUserModel(rowUser)

	if !user.HasPermision() {
		return enerr.E("No permission to start game")
	}

	// TODO:
	engine, err := uc.RoomRepo.GetEngine(roomName)

	engine.Exit()

	err = tx.Commit()
	if err != nil {
		return enerr.E(op, err, enerr.Internal)
	}

	return nil

}
