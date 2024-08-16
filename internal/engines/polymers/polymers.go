package polymers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
)

var (
	ErrAlreadyStarted = errors.New("engine was started")
	ErrMaxPlayers     = errors.New("all players places are filled")
)

func New(log *slog.Logger, cfg PolymersEngineConfig) *PolymersEngine {
	// var src cryptorand.CryptoSource
	src := rand.NewSource(time.Now().UnixNano())

	eng := &PolymersEngine{
		log:        log.With(slog.String("engine", "PolymersEngine")),
		Bag:        NewGameBag(cfg.Elements),
		ActionChan: make(chan models.Action),
		handlers:   map[string]HandlerFunc{},
		unicast:    cfg.Unicast,
		broadcast:  cfg.Broadcast,
		timerInt:   cfg.TimerInt,
		MaxPlayers: cfg.MaxPlayers,
		rnd:        rand.New(src),
	}
	eng.StateMachine = stateMachine{
		Current: OBTAIN,
		States: map[stateInt]State{
			OBTAIN: {
				handlers: map[string]HandlerFunc{
					"GetElement": eng.GetElement(),
				},
			},
		},
	}

	return eng
}

type PolymersEngineConfig struct {
	Elements   map[string]int
	TimerInt   int
	Unicast    unicastFunction
	Broadcast  broadcastFunction
	MaxPlayers int
}

type PolymersEngine struct {
	log          *slog.Logger
	started      bool
	StateMachine stateMachine
	Bag          GameBag
	// Канал для обработки действий игроков
	ActionChan chan models.Action
	handlers   map[string]HandlerFunc

	players    []string
	MaxPlayers int

	unicast   unicastFunction
	broadcast broadcastFunction

	timerInt int
	mu       sync.Mutex
	rnd      *rand.Rand
}

// unicastFunction accepts first argument userID
type unicastFunction func(string, common.Message)
type broadcastFunction func(common.Message)

// TODO: OnlyOnce запустить только раз.
func (engine *PolymersEngine) Start() {
	if engine.started {
		// TODO:
		return
	}
	engine.started = true
	go func() {
		for {
			select {
			case e := <-engine.ActionChan:
				engine.mu.Lock()
				state, err := engine.StateMachine.States[engine.StateMachine.Current].Input(e)
				if err != nil {
					engine.log.Error("error while handling action with state", sl.Err(err), slog.String("state", engine.StateMachine.Current.String()))
				}
				// Changing state if needed
				if state > NO_TRANSITION {
					engine.StateMachine.Current = state
				}
				engine.mu.Unlock()
			}
			engine.log.Debug("Engine selected action")
		}
	}()
	engine.log.Debug("Broadcast for starting engine")
	engine.broadcast(common.Message{Type: common.HUB_STARTGAME, Ok: true})
}

func (engine *PolymersEngine) Input(e models.Action) {
	engine.ActionChan <- e
}
func (engine *PolymersEngine) PreHook() map[string]any {

	return map[string]any{"engine": engine}
}

func (engine *PolymersEngine) AddPlayer(name string) error {
	if engine.started {
		return ErrAlreadyStarted
	}
	if len(engine.players) >= engine.MaxPlayers {
		return ErrMaxPlayers
	}
	engine.players = append(engine.players, name)
	return nil
}

func (engine *PolymersEngine) RemovePlayer(name string) error {
	if engine.started {
		return ErrAlreadyStarted
	}
	engine.players = removeByValue(engine.players, name)
	return nil
}

func (engine *PolymersEngine) MarshalJSON() ([]byte, error) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	eng := struct {
		Started bool
		State   string
		Bag     GameBag
		Players []string
	}{
		engine.started,
		engine.StateMachine.Current.String(),
		engine.Bag,
		engine.players,
	}
	return json.Marshal(eng)
}
