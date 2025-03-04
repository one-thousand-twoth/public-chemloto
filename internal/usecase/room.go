package usecase

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type CreateRoomRequest struct {
	Name     string `json:"name" validate:"required,min=1,safeinput"`
	EngineID string
}
type Response struct {
	Rooms any      `json:"rooms"`
	Error []string `json:"error"`
}

func CreateRoom(req CreateRoomRequest, db *sql.DB) (*entities.Room, error) {

	const op = "server.handlers.CreateRoom"

	params := database.InsertRoomParams{
		Name:   req.Name,
		Engine: req.EngineID,
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
	// if err := s.hub.AddNewRoom(req.CreateRoomRequest); err != nil {
	// 	log.Error("failed to add room", sl.Err(err))
	// 	switch validateErr := err.(type) {
	// 	case validator.ValidationErrors:
	// 		encode(w, r, http.StatusBadRequest, Response{Error: appvalidation.ValidationError(validateErr)})
	// 		return
	// 	}
	// 	if enerr.KindIs(enerr.Exist, err) {
	// 		encode(w, r, http.StatusConflict, Response{Error: []string{"Комната уже существует"}})
	// 		return
	// 	}
	// 	encode(w, r, http.StatusConflict, Response{Error: []string{"Сервер не смог создать комнату"}})
	// 	return
	// }
	// s.log.Info("Room created", "name", req.Name, "time", req.Time)
	// s.hub.SendMessageOverChannel("default", models.Message{Type: websocket.TextMessage, Body: []byte(req.Name)})
	// encode(w, r, http.StatusOK, Response{Rooms: s.hub.Rooms, Error: []string{}})
	room := &entities.Room{
		Name:   row.Name,
		Engine: entities.ExternalEngine{},
	}
	return room, nil
}
