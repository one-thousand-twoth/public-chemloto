package usecase

import (
	"errors"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
)

func AddRegularChannel(repo *repository.ChannelsRepository, name string, fn entities.InitFunction) (*entities.Channel, error) {
	if fn == nil {
		return nil, enerr.E(errors.New("nil function"))
	}
	row, err := repo.AddRegularChannel(name, fn)

	return row, err

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

type ChannelRepository interface {
	Get(channel string) ([]string, bool)
	Add(channel string, connection string)
	Remove(channel string, connection string)
	GetChannelFunc(channel string) (func(chan common.Message), bool)
	SetChannelFunc(channel string, fun func(chan common.Message))
}

func SubscribeToRoom(
	channelsRepo *repository.ChannelsRepository,
	roomRepo *repository.RoomRepository,
	roomName string,
	user *entities.User,
) error {
	const op enerr.Op = "usecase.subscribtions/SubscribeToRoom"

	// if data.Target == "" || data.Name == "" {
	// 	return enerr.E(op, "empty field", enerr.Validation)
	// }

	err := roomRepo.SubscribeToRoom(roomName, user)
	if err != nil {
		return err
	}

	return nil

}
