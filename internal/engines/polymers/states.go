package polymers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	TRADE
	COMPLETED
)

type stateInt int

// Реализация FSM (Finite State Machine)
type stateMachine struct {
	Current stateInt
	States  map[stateInt]State
}

type State interface {
	Handle(e models.Action, player *Participant) (stateInt, error)
	PreHook()
	// Handlers() map[string]HandlerFunc
	Update() (stateInt, error)
	MarshalJSON() ([]byte, error)
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
func (s SimpleState) Update() (stateInt, error) {
	return NO_TRANSITION, nil
}
func (s SimpleState) Handlers() map[string]HandlerFunc {
	return s.handlers
}

// MarshalJSON - метод, который всегда возвращает пустой JSON-объект
func (s SimpleState) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
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
	// fmt.Println("find handler for action")
	st = handle(e)
	return st, nil
}
func (s SimpleState) PreHook() {

}

type ObtainState struct {
	ticker  *time.Ticker
	maxDur  int
	currDur time.Duration
	step    int
	done    chan struct{}
	eng     *PolymersEngine
	SimpleState
	startTime time.Time
}

func (s *ObtainState) Handlers() map[string]HandlerFunc {
	return s.handlers
}

// MarshalJSON - метод, который всегда возвращает пустой JSON-объект
func (s *ObtainState) MarshalJSON() ([]byte, error) {
	st := struct {
		Timer int
	}{
		int((s.currDur * time.Second).Seconds() - time.Since(s.startTime).Seconds()),
	}
	return json.Marshal(st)
}
func (s *ObtainState) PreHook() {
	// empty action, bc no needed
	s.handlers["GetElement"](models.Action{})
	s.eng.log.Debug("Reseting Timer")
	s.ticker.Reset(s.currDur * time.Second)
	s.startTime = time.Now()
}

// GetElement now
func (s *ObtainState) Update() (stateInt, error) {
	s.eng.log.Debug("Updating ObtainState")
	st := s.handlers["GetElement"](models.Action{})
	if st == OBTAIN || st == NO_TRANSITION {
		s.currDur = s.incrementalTime()
		s.step += 1
		s.eng.log.Debug("new Timer", "currDur", s.currDur, "step", s.step)
		s.ticker.Reset(s.currDur * time.Second)
		s.eng.broadcast(common.Message{
			Type: common.ENGINE_ACTION,
			Ok:   true,
			Body: map[string]any{"Action": "NewTimer", "Value": (s.currDur * time.Second).Seconds()},
		})
		s.startTime = time.Now()
	}

	return st, nil
}
func (s *ObtainState) incrementalTime() time.Duration {
	if s.step <= 10 {
		// Для первых 10 ходов вычисляем значение как логарифм по основанию 2
		// Добавляем 1 чтобы избежать логарифма от 0
		dur := time.Duration(math.Log2(float64(s.step)+1) + 1)
		if dur < time.Duration(s.maxDur)*time.Second {
			return dur
		}
	}
	// Для всех остальных ходов возвращаем заданное время
	return time.Duration(s.maxDur)

}

func (eng *PolymersEngine) NewObtainState(timerInt int) (state *ObtainState) {
	state = &ObtainState{
		eng:     eng,
		done:    make(chan struct{}),
		maxDur:  timerInt,
		currDur: time.Duration(1),
		step:    0,
		// Cпециально недостижимое число, в дальнейшем функции будут переопределять значение
		ticker:      time.NewTicker(time.Hour * 100),
		SimpleState: NewState(),
	}
	state.Add("GetElement", GetElement(eng), true)
	state.Add("RaiseHand", RaiseHand(eng), false)
	go func() {
		for t := range state.ticker.C {
			eng.Internal <- t
		}
		eng.log.Error("exiting timer loop")
	}()
	return
}
