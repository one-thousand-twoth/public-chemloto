package repository

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database/stores"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepository struct {
	messageChannels *stores.StreamStore
	db              *sql.DB
	queries         *database.Queries
}

func NewUserRepo(db *sql.DB, memMessageStream *stores.StreamStore) *UserRepository {
	return &UserRepository{
		messageChannels: memMessageStream,
		db:              db,
		queries:         database.New(db),
	}
}

func (repo *UserRepository) CreateUser(params database.InsertUserParams) (*entities.User, error) {
	const op enerr.Op = "repository.user/PatchUser"
	tx, err := repo.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
	}
	defer tx.Rollback()
	queries := repo.queries.WithTx(tx)
	row, err := queries.InsertUser(context.TODO(), params)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return nil, enerr.E(op, "Пользователь с таким именем уже существует", enerr.Exist)
			}
		}
		// encode(w, r, http.StatusConflict, Response{Error: []string{"Пользователь с таким именем уже существует"}})
		return nil, enerr.E(op, err, enerr.Internal)
	}
	// NOTE: ВЫЗЫВАЕТ ОШИБКУ FOREIGN KEY если нет такого канала
	// _, err = queries.SubscribeToGroup(context.TODO(),
	// 	database.SubscribeToGroupParams{
	// 		ChannelID: 1,
	// 		UserID:    0,
	// 	})
	// if err != nil {
	// 	return nil, enerr.E(op, err, enerr.Database)
	// }

	// TODO:
	repo.messageChannels.Add(row.Name, make(entities.MessageStream, 1000))

	if err := tx.Commit(); err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
	}
	user := &entities.User{
		ID:          entities.ID(row.ID),
		Name:        row.Name,
		Apikey:      row.Apikey,
		Room:        row.Room.String,
		Role:        common.Role(row.Role),
		MessageChan: repo.messageChannels.Get(row.Name),
	}
	return user, nil
}

func (repo *UserRepository) GetUserByApikey(apikey string) (*entities.User, error) {
	row, err := repo.queries.GetUserByApikey(context.TODO(), apikey)
	if err != nil {
		return nil, err
	}

	user := entities.User{
		ID:          entities.ID(row.ID),
		Name:        row.Name,
		Apikey:      row.Apikey,
		Room:        row.Room.String,
		Role:        common.Role(row.Role),
		MessageChan: repo.messageChannels.Get(row.Name),
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByID(id entities.ID) (*entities.User, error) {
	row, err := repo.queries.GetUserByID(context.TODO(), int64(id))
	if err != nil {
		return nil, err
	}

	user := entities.User{
		ID:          entities.ID(row.ID),
		Name:        row.Name,
		Apikey:      row.Apikey,
		Room:        row.Room.String,
		Role:        common.Role(row.Role),
		MessageChan: repo.messageChannels.Get(row.Name),
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByName(name string) (*entities.User, error) {
	const op enerr.Op = "repository.user/GetUserByName"
	row, err := repo.queries.GetUserByName(context.TODO(), name)
	if err != nil {
		return nil, enerr.E(op, err)
	}

	user := entities.User{
		ID:          entities.ID(row.ID),
		Name:        row.Name,
		Apikey:      row.Apikey,
		Room:        row.Room.String,
		Role:        common.Role(row.Role),
		MessageChan: repo.messageChannels.Get(row.Name),
	}

	return &user, nil
}

func (repo *UserRepository) GetUsers() ([]entities.User, error) {
	const op enerr.Op = "repository.user/GetUsers"
	userRows, err := repo.queries.GetUsers(context.TODO())
	if err != nil {
		return nil, enerr.E(op, err)
	}
	users := make([]entities.User, 0, len(userRows))
	for _, row := range userRows {
		users = append(users, entities.User{
			ID:          entities.ID(row.ID),
			Name:        row.Name,
			Apikey:      row.Apikey,
			Room:        row.Room.String,
			Role:        common.Role(row.Role),
			MessageChan: repo.messageChannels.Get(row.Name),
		})
	}
	return users, nil
}

func (repo *UserRepository) GetRoomSubscribers(roomname string) ([]entities.User, error) {
	rows, err := repo.queries.GetUsersByRoom(context.TODO(),
		sql.NullString{
			String: roomname,
			Valid:  true,
		})
	if err != nil {
		return nil, err
	}
	users := make([]entities.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, entities.User{
			ID:          entities.ID(row.ID),
			Name:        row.Name,
			Apikey:      row.Apikey,
			Room:        row.Room.String,
			Role:        common.Role(row.Role),
			MessageChan: repo.messageChannels.Get(row.Name),
		})
	}
	return users, nil
}
