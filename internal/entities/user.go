package entities

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
)

type User struct {
	Name     string
	Apikey   string
	Room     string
	Role     common.Role
	conn     string
	channels []string
}

func NewUser(name string, apikey string, conn string, role common.Role, channels []string) *User {
	return &User{
		Name:     name,
		Apikey:   apikey,
		conn:     conn,
		channels: channels,
		Role:     role,
	}
}
