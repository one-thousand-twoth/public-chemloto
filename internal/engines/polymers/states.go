package polymers

import (
	"errors"
	"fmt"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
)

var (
	ErrNoHandler    = errors.New("no handler")
	ErrNoAuthorized = errors.New("no authorized player")
)

//go:generate stringer -type=stateInt
const (
	NO_TRANSITION stateInt = iota
	OBTAIN
	HAND
	COMPLETED
)

type stateInt int

type stateMachine struct {
	Current stateInt
	States  map[stateInt]State
}

type State struct {
	handlers map[string]HandlerFunc
	secure   map[string]bool
}

func NewState() State {
	return State{
		handlers: make(map[string]HandlerFunc),
		secure:   make(map[string]bool),
	}
}

func (s State) Add(state string, fun HandlerFunc, secure bool) State {
	s.handlers[state] = fun
	s.secure[state] = secure
	return s
}

// Handle can return ErrNoAuthorized if action cannot be sent by player.
// ErrNoHandler if no handler was found for action
func (s State) Handle(e models.Action, player models.Player) (stateInt, error) {
	// Default not to transition to another state
	st := NO_TRANSITION

	action, ok := e.Envelope["Action"].(string)
	if !ok {
		return st, errors.New("failed to extract Action field")
	}
	if s.secure[action] {
		if player.Role < common.Judge_Role {
			return NO_TRANSITION, fmt.Errorf("%w for %s action", ErrNoAuthorized, action)
		}
	}
	handle, ok := s.handlers[action]
	if !ok {
		return st, fmt.Errorf("%w for %s action", ErrNoHandler, action)
	}
	fmt.Println("find handler for action")
	st = handle(e)
	return st, nil
}
