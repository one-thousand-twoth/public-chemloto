package usecase

import (
	"errors"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
)

func AddRegularChannel(repo *repository.ChannelsRepository, name string, fn func()) (entities.ID, error) {
	if fn == nil {
		return -1, enerr.E(errors.New("nil function"))
	}
	id, err := repo.AddRegularChannel(name, fn)

	return id, err
}

func GetRegularChannel(repo *repository.ChannelsRepository, id entities.ID) error {
	err := repo.GetChannelByID(id)
	if err != nil {
		return err
	}
	return nil
}

func SubscribeToChannel(repo *repository.ChannelsRepository, id entities.ID, user entities.User) error {

	err := repo.SubscribeTo(int64(id), user)
	if err != nil {
		return err
	}
	return nil

}
