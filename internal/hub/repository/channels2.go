package repository

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
)

type ChannelsRepository struct {
	channelFunctions map[string]entities.InitFunction

	queries *database.Queries
}

func NewChannelsRepo(db *sql.DB) *ChannelsRepository {
	return &ChannelsRepository{
		channelFunctions: map[string]entities.InitFunction{
			"default": func(chan common.Message) {
				//do not do anything at that moment
			},
		},
		queries: database.New(db),
	}
}

func (repo *ChannelsRepository) AddRegularChannel(name string, fn entities.InitFunction) (*entities.Channel, error) {
	row, err := repo.queries.InsertRegularChannel(context.TODO(), name)
	if err != nil {
		return nil, err
	}
	repo.channelFunctions[name] = fn
	channel := &entities.Channel{
		ID:   entities.ID(row.ID),
		Name: row.Name,
		Type: row.Type,
		Fn:   fn,
	}
	return channel, nil
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
