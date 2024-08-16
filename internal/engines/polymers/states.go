package polymers

import (
	"errors"
	"fmt"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
)

var (
	ErrNoHandler = errors.New("no handler")
)

//go:generate stringer -type=stateInt
const (
	NO_TRANSITION stateInt = iota
	OBTAIN
	HANDLE
	COMPLETED
)

type stateInt int

type stateMachine struct {
	Current stateInt
	States  map[stateInt]State
}

type State struct {
	handlers map[string]HandlerFunc
}

func (s State) Input(e models.Action) (stateInt, error) {
	// Default not to transition to another state
	st := NO_TRANSITION
	action, ok := e.Envelope["Action"].(string)
	if !ok {
		return st, errors.New("failed to extract Action field")
	}
	handle, ok := s.handlers[action]
	if !ok {
		return st, fmt.Errorf("%w for %s action", ErrNoHandler, action)
	}
	st = handle(e)
	return st, nil
}
