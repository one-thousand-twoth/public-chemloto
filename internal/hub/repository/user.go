package repository

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type messageChan chan common.Message

type UserRepository struct {
	messageChannels map[string]messageChan
	db              *sql.DB
	queries         *database.Queries
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		messageChannels: map[string]messageChan{},
		db:              db,
		queries:         database.New(db),
	}
}

func (repo *UserRepository) CreateUser(params database.InsertUserParams) (*entities.User, error) {
	const op enerr.Op = "usecase.user/PatchUser"
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
				return nil, enerr.E(op, err, enerr.Exist)
			}
		}
		// encode(w, r, http.StatusConflict, Response{Error: []string{"Пользователь с таким именем уже существует"}})
		return nil, enerr.E(op, err, enerr.Internal)
	}
	// NOTE: ВЫЗЫВАЕТ ОШИБКУ FOREIGN KEY если нет такого канала
	_, err = queries.InsertChannelSubscribeByChannelName(context.TODO(),
		database.InsertChannelSubscribeByChannelNameParams{
			Name:   "default",
			UserID: row.ID,
		})
	if err != nil {
		return nil, enerr.E(op, err, enerr.Database)
	}

	repo.messageChannels[row.Name] = make(messageChan)

	user := &entities.User{
		ID:          entities.ID(row.ID),
		Name:        row.Name,
		Apikey:      row.Apikey,
		Room:        row.Room.String,
		Role:        common.Role(row.Role),
		MessageChan: repo.messageChannels[row.Name],
	}

	if err := tx.Commit(); err != nil {
		return nil, enerr.E(op, err, enerr.Internal)
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
		MessageChan: make(chan common.Message),
	}

	return &user, nil
}
