package models

import (
	"errors"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
)

type Action struct {
	Player   string
	Envelope map[string]any
}

type Player struct {
	Name string
	Role common.Role
}

var (
	ErrAlreadyStarted = errors.New("engine was started")
	ErrMaxPlayers     = errors.New("all players places are filled")
)
