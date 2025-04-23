package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database/stores"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type RoomRepository struct {
	engines *stores.EnginesStore
	streams *stores.StreamStore
	db      *sql.DB
	queries *database.Queries
}

func NewRoomRepo(db *sql.DB, memMessageStream *stores.StreamStore) *RoomRepository {
	return &RoomRepository{
		db:      db,
		queries: database.New(db),
		engines: stores.NewEngineStore(),
		streams: memMessageStream,
	}
}

// Предполгается что вы точно знаете, что комната была создана и вы вызываете ее в блокирующей транзакции
//
// Эта функция опасна, так как может неявно отдать нулевое значение
// Но я закрываю глаза, так как она используется лишь в одном случае, из-за плохой моей архитектуры :Э
func (repo *RoomRepository) GetEngine(name string) (models.Engine, error) {
	const op enerr.Op = "repository.room/GetEngine"

	return repo.engines.Get(name), nil
}

func (repo *RoomRepository) AddRoom(name string, engine models.Engine) (*entities.Room, error) {
	const op enerr.Op = "repository.room/AddRoom"
	tx, err := repo.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
	}
	defer tx.Rollback()
	queries := repo.queries.WithTx(tx)

	rowRoom, err := queries.InsertRoom(context.TODO(), database.InsertRoomParams{
		Name:   name,
		Engine: name,
	})
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return nil, enerr.E(op, "Комната уже существует", enerr.Exist)
			}
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
				return nil, enerr.E(op, "Комната уже существует", enerr.Exist)
			}
		}
		return nil, enerr.E(op, err, enerr.Internal)
	}

	repo.engines.Add(rowRoom.Name, engine)

	if err := tx.Commit(); err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
	}

	room := &entities.Room{
		Name:   rowRoom.Name,
		Engine: engine,
	}
	return room, nil
}

func (repo *RoomRepository) GetRooms() ([]*entities.Room, error) {
	const op enerr.Op = "repository.room/GetRooms"
	rows, err := repo.queries.GetRooms(context.TODO())
	if err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
	}

	rooms := make([]*entities.Room, 0, len(rows))
	for _, v := range rows {
		rooms = append(rooms, &entities.Room{
			Name:   v.Name,
			Engine: repo.engines.Get(v.Name),
		})
	}
	return rooms, nil
}

func (repo *RoomRepository) GetRoom(name string) (*entities.Room, error) {
	const op enerr.Op = "repository.room/GetRooms"
	row, err := repo.queries.GetRoomByName(context.TODO(), name)
	if err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
	}

	room := &entities.Room{
		Name:   row.Name,
		Engine: repo.engines.Get(row.Name),
	}
	return room, nil
}
func (repo *RoomRepository) SubscribeToRoom(name string, user *entities.User) error {
	const op enerr.Op = "repository.room/SubscribeToRoom"
	tx, err := repo.db.BeginTx(context.TODO(), nil)
	defer tx.Rollback()
	if err != nil {
		return enerr.E(op, err, enerr.Database)
	}
	queries := repo.queries.WithTx(tx)

	room, err := queries.GetRoomByName(context.TODO(), name)
	if err != nil {
		return enerr.E(op, err, enerr.Database)
	}

	err = queries.UpdateUserRoom(context.TODO(), database.UpdateUserRoomParams{
		Room: sql.NullString{String: name, Valid: true},
		ID:   int64(user.ID),
	})
	if err != nil {
		return enerr.E(op, err)
	}
	engine := repo.engines.Get(room.Name)
	err = engine.AddParticipant(models.Participant{
		Name: user.Name,
		Role: user.Role,
	})
	if err != nil {
		return enerr.E(op, err)
	}

	if err := tx.Commit(); err != nil {
		return enerr.E(op, err, enerr.Internal)
	}

	return nil
}

func (repo *RoomRepository) GetRoomUsers(ctx context.Context, roomname string) ([]*entities.User, error) {
	const op enerr.Op = "repository.room/GetRoomUsers"

	rows, err := repo.queries.GetUsersByRoom(ctx, sql.NullString{
		String: roomname,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rows = make([]database.User, 0)

		} else {
			return nil, err
		}
	}
	users := make([]*entities.User, 0, len(rows))
	for _, v := range rows {
		usM := entities.ToUserModel(v)
		usM.MessageChan = repo.streams.Get(v.Name)
		users = append(users, &usM)
	}
	return users, nil
}
