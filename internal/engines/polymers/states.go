package polymers

import (
	"errors"
	"fmt"
	"time"

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

type State interface {
	Handle(e models.Action, player *Participant) (stateInt, error)
	PreHook()
	// Handlers() map[string]HandlerFunc
	Update()
}

// SimpleState
// Use it for state with no need for struct model
type SimpleState struct {
	handlers map[string]HandlerFunc
	secure   map[string]bool
}

func NewState() SimpleState {
	return SimpleState{
		handlers: make(map[string]HandlerFunc),
		secure:   make(map[string]bool),
	}
}

func (s SimpleState) Handlers() map[string]HandlerFunc {
	return s.handlers
}

// Add will add Handler for State, {secure} defines handler that works only for Judge or Admin
func (s SimpleState) Add(action string, fun HandlerFunc, secure bool) SimpleState {
	s.handlers[action] = fun
	s.secure[action] = secure
	return s
}

// Handle can return ErrNoAuthorized if action cannot be sent by player.
// ErrNoHandler if no handler was found for action
func (s SimpleState) Handle(e models.Action, player *Participant) (stateInt, error) {
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
func (s SimpleState) PreHook() {

}

type ObtainState struct {
	ticker   *time.Ticker
	done     chan struct{}
	eng      *PolymersEngine
	handlers map[string]HandlerFunc
	secure   map[string]bool
}

func (s *ObtainState) Add(action string, fun HandlerFunc, secure bool) *ObtainState {
	s.handlers[action] = fun
	s.secure[action] = secure
	return s
}
func (s *ObtainState) Handlers() map[string]HandlerFunc {
	return s.handlers
}

func (s *ObtainState) Handle(e models.Action, player *Participant) (stateInt, error) {
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

func (s *ObtainState) PreHook() {
	// empty action, bc no needed
	s.handlers["GetElement"](models.Action{})
	s.ticker.Reset(time.Second * 2)
}

// GetElement now
func (s *ObtainState) Update() {

}

func (eng *PolymersEngine) NewObtainState() (state *ObtainState) {
	state = &ObtainState{
		eng:      eng,
		done:     make(chan struct{}),
		ticker:   time.NewTicker(time.Second),
		handlers: make(map[string]HandlerFunc),
		secure:   make(map[string]bool),
	}
	state.Add("GetElement", GetElement(eng), true)
	state.Add("RaiseHand", RaiseHand(eng), false)
	go func() {
		for t := range eng.ticker.C {
			eng.Internal <- t
		}
	}()
	return
}
