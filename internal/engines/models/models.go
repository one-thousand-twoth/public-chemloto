package models

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
)

type Engine interface {
	// Получить текущее состояние, например при перезагрузке страницы
	PreHook() map[string]any
	// Обработать событие
	Input(Action)
	Start()
	GetResults() [][]string
	AddParticipant(Participant) error
	RemoveParticipant(name string) error
	Exit()
}

type EngineStatus int

//go:generate stringer -type=EngineStatus
const (
	STATUS_WAITING EngineStatus = iota
	STATUS_STARTED
	STATUS_COMPLETED
)

// UnicastFunction accepts first argument userID
type UnicastFunction func(userId string, msg common.Message)
type BroadcastFunction func(common.Message)

type Action struct {
	Player   string
	Envelope map[string]any
}

type Participant struct {
	Name string
	Role common.Role
}
