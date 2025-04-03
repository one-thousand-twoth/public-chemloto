package entities

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
)

type ID int64

type User struct {
	ID          ID                  `json:"id"`
	Name        string              `json:"username"`
	Apikey      string              `json:"-"`
	Room        string              `json:"room"`
	Role        common.Role         `json:"role"`
	MessageChan chan common.Message `json:"-"`
	// channels    []string            `json:"channels"`
}

func (u *User) SubscribeToRoom(roomName string) error {
	if u.Room != "" {
		return enerr.E("user already subscribed to a room")
	}

	u.Room = roomName
	return nil

}

func NewUser(name string, apikey string, conn string, role common.Role, channels []string) *User {
	return &User{
		Name:   name,
		Apikey: apikey,
		// channels: channels,
		Role: role,
	}
}
