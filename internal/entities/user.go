package entities

import (
	"encoding/json"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
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

func ToUserModel(u database.User) User {

	return User{
		ID:     ID(u.ID),
		Name:   u.Name,
		Apikey: u.Apikey,
		Room:   u.Room.String,
		Role:   common.Role(u.Role),
		// MessageChan: nil,
	}
}

func (r *User) MarshalJSON() ([]byte, error) {
	user := struct {
		Name string `json:"username"`
		Room string `json:"room"`
		Role string `json:"role"`
	}{r.Name, r.Room, r.Role.String()}
	return json.Marshal(user)
}

func (u *User) HasPermision() bool {
	if u.Role < common.Judge_Role {
		return false
	}
	return true
}
func (u *User) IsInRoom() bool {
	if u.Room != "" {
		return true
	}
	return false
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
