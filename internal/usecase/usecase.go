package usecase

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
)

type RoomRepository interface {
	AddRoom(name string, engine models.Engine) (*entities.Room, error)
	GetRooms() ([]*entities.Room, error)
	GetRoom(name string) (*entities.Room, error)
	SubscribeToRoom(name string, user *entities.User) error
}

type DBInterface interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Usecases struct {
	RoomRepo  RoomRepository
	UserRepo  *repository.UserRepository
	GroupRepo *repository.GroupsRepository

	db      *sql.DB
	queries database.Queries
}

func NewUsecase(db *sql.DB) *Usecases {

	return &Usecases{
		RoomRepo:  repository.NewRoomRepo(db),
		UserRepo:  repository.NewUserRepo(db),
		GroupRepo: repository.NewGroupsRepo(db),
		db:        db,
		queries:   *database.New(db),
	}

}
