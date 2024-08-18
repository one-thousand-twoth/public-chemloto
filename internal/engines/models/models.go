package models

import "github.com/anrew1002/Tournament-ChemLoto/internal/common"

type Action struct {
	Player   string
	Envelope map[string]any
}

type Player struct {
	Name       string
	Role       common.Role
	RaisedHand bool
}
