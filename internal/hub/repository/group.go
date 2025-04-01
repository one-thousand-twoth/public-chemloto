package repository

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
)

type GroupsRepository struct {
	groupFunctions map[string]entities.InitFunction

	queries *database.Queries
}

func NewGroupsRepo(db *sql.DB) *GroupsRepository {
	return &GroupsRepository{
		groupFunctions: map[string]entities.InitFunction{
			"default": func(chan common.Message) {
				//do not do anything at that moment
			},
		},
		queries: database.New(db),
	}
}

func (repo *GroupsRepository) AddRegularGroup(name string, fn entities.InitFunction) (*entities.Group, error) {
	row, err := repo.queries.InsertGroup(context.TODO(), name)
	if err != nil {
		return nil, err
	}
	repo.groupFunctions[name] = fn
	channel := &entities.Group{
		ID:   entities.ID(row.ID),
		Name: row.Name,
		Fn:   fn,
	}
	return channel, nil
}
func (repo *GroupsRepository) GetGroupByID(id entities.ID) error {
	_, err := repo.queries.GetGroupByID(context.TODO(), int64(id))
	if err != nil {
		return err
	}
	return nil
}

func (repo *GroupsRepository) GetAllUserGroups(userID entities.ID) ([]*entities.Group, error) {
	rows, err := repo.queries.GetGroupByUserID(context.TODO(), int64(userID))
	if err != nil {
		return nil, err
	}
	channels := make([]*entities.Group, 0, len(rows))
	// for _, v := range rows {
	// 	fn, ok := repo.groupFunctions[v.RoomName.String]
	// 	if !ok {
	// 		fn = func(chan common.Message) {}
	// 	}
	// 	channels = append(channels, &entities.Group{
	// 		ID:   entities.ID(v.ID),
	// 		Name: v.Name,
	// 		Type: v.Type,
	// 		Fn:   fn,
	// 	})
	// }
	return channels, err
}

func (repo *GroupsRepository) SubscribeTo(channelID int64, user entities.User) error {
	params := database.SubscribeToGroupParams{
		ChannelID: channelID,
		UserID:    int64(user.ID),
	}
	err := repo.queries.SubscribeToGroup(context.TODO(), params)
	if err != nil {
		return err
	}
	return nil
}
