package models

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
)

// UnicastFunction accepts first argument userID
type UnicastFunction func(userId string, msg common.Message)
type BroadcastFunction func(common.Message)

type Action struct {
	Player   string
	Envelope map[string]any
}

type Player struct {
	Name string
	Role common.Role
}
