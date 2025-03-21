package usecase

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

func CreateRoom(req CreateRoomRequest, log *slog.Logger, db *sql.DB) (*entities.Room, error) {

	const op = "server.handlers.CreateRoom"

	engines.NewEngine(req.Type, req.Name, log, req.EngineConfig)

	params := database.InsertRoomParams{
		Name:   req.Name,
		Engine: req.Name,
	}
	row, err := database.New(db).InsertRoom(context.TODO(), params)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return nil, enerr.E(op, err, enerr.Exist)
			}
		}
		return nil, enerr.E(op, err, enerr.Internal)
	}

	// TODO: Create channel for room

	room := &entities.Room{
		Name:   row.Name,
		Engine: entities.ExternalEngine{},
	}
	return room, nil
}

func GetRooms(db *sql.DB) []entities.Room {

	rows, err := database.New(db).GetRooms(context.TODO())
	if err != nil {
		panic(err)
	}

	rooms := make([]entities.Room, 0, len(rows))
	for _, v := range rows {
		rooms = append(rooms, entities.Room{
			Name:   v.Name,
			Engine: entities.ExternalEngine{},
		})
	}

	return rooms
}
