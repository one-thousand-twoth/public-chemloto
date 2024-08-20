package common

import "encoding/json"

//  at that time - deprecated
type User struct {
	// Id       string
	Username string
	Score    int
	Room     string
	// Password string
	Admin bool
}

//  at that time - deprecated
type Room struct {
	Name       string         `json:"roomName" validate:"required"`
	Time       int            `json:"time" validate:"required_if=isAuto true,gt=0"`
	Max_partic int            `json:"maxPlayers" validate:"required"`
	Elements   map[string]int `json:"elementCounts" validate:"required"`
	IsAuto     bool           `json:"isAuto"`
}

type Role int

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

//go:generate stringer -type=Role
const (
	NONE Role = iota
	Player_Role
	Judge_Role
	Admin_Role
)
