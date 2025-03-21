package repository

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
)

type ChannelsRepository struct {
	channelFunctions map[string]func()

	queries *database.Queries
}

func NewChannelsRepo(db *sql.DB) *ChannelsRepository {
	return &ChannelsRepository{
		channelFunctions: map[string]func(){},
		queries:          database.New(db),
	}
}

func (repo *ChannelsRepository) AddRegularChannel(name string, fn func()) (entities.ID, error) {
	row, err := repo.queries.InsertRegularChannel(context.TODO(), name)
	if err != nil {
		return -1, err
	}
	repo.channelFunctions[name] = fn
	return entities.ID(row.ID), nil
}
func (repo *ChannelsRepository) GetChannelByID(id entities.ID) error {
	_, err := repo.queries.GetChannelByID(context.TODO(), int64(id))
	if err != nil {
		return err
	}
	return nil
}

func (repo *ChannelsRepository) SubscribeTo(channelID int64, user entities.User) error {
	params := database.InsertChannelSubscribeParams{
		ChannelID: channelID,
		UserID:    int64(user.ID),
	}
	_, err := repo.queries.InsertChannelSubscribe(context.TODO(), params)
	if err != nil {
		return err
	}
	return nil
}
