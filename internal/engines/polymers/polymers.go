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
		Checks:     cfg.Checks,
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
			OBTAIN: NewState().
				Add("GetElement", eng.GetElement(), true).
				Add("RaiseHand", eng.RaiseHand(), false),
			HAND: NewState().
				Add("RaiseHand", eng.RaiseHand(), false).
				Add("Check", eng.Check(), true),
		},
	}

	return eng
}

type PolymersEngineConfig struct {
	Elements   map[string]int
	Checks     Checks
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
	Checks     Checks

	players    []models.Player
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
				engine.log.Debug("locked engine")
				func() {
					player, err := engine.GetPlayer(e.Player)
					if err != nil {
						engine.log.Error("Unknown player", slog.String("name", e.Player))
						return
					}
					state, err := engine.StateMachine.States[engine.StateMachine.Current].Handle(e, player)
					if err != nil {
						if errors.Is(err, ErrNoAuthorized) {
							engine.log.Error("no authorized")
							engine.unicast(e.Player, common.Message{Type: common.UNDEFINED, Ok: false, Errors: []string{"Недостаточно прав"}})
						}
						engine.log.Error("error while handling action with state", sl.Err(err), slog.String("state", engine.StateMachine.Current.String()))

						return
					}
					// Changing state if needed
					if state > NO_TRANSITION {
						engine.log.Info("Changing game state",
							slog.String("new state", state.String()),
							slog.String("old state", engine.StateMachine.Current.String()))
						engine.StateMachine.Current = state
					}
				}()
				engine.mu.Unlock()
				engine.log.Debug("unlocked engine")
			}
			engine.log.Debug("Engine selected action")
		}
	}()
	engine.log.Debug("Broadcast for starting engine")
	engine.broadcast(common.Message{Type: common.HUB_STARTGAME, Ok: true})
}

func (engine *PolymersEngine) Input(e models.Action) {
	if !engine.started {
		engine.unicast(e.Player, common.Message{
			Type:   common.UNDEFINED,
			Ok:     false,
			Errors: []string{"Игра ещё не начата"},
		})
		return
	}
	engine.ActionChan <- e
}
func (engine *PolymersEngine) PreHook() map[string]any {
	return map[string]any{"engine": engine}
}

func (engine *PolymersEngine) AddPlayer(player models.Player) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	if engine.started {
		return ErrAlreadyStarted
	}
	if len(engine.players) >= engine.MaxPlayers {
		return ErrMaxPlayers
	}
	engine.players = append(engine.players, player)
	return nil
}

func (engine *PolymersEngine) RemovePlayer(name string) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	if engine.started {
		return ErrAlreadyStarted
	}
	for i := 0; i < len(engine.players); i++ {
		if engine.players[i].Name == name {
			engine.players = append(engine.players[:i], engine.players[i+1:]...)
			break
		}
	}
	return nil
}
func (engine *PolymersEngine) GetPlayer(name string) (models.Player, error) {
	for i := 0; i < len(engine.players); i++ {
		if engine.players[i].Name == name {
			return engine.players[i], nil
		}
	}
	return models.Player{}, errors.New("unknown player")
}

func (engine *PolymersEngine) MarshalJSON() ([]byte, error) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	eng := struct {
		Started bool
		State   string
		Bag     GameBag
		Players []models.Player
	}{
		engine.started,
		engine.StateMachine.Current.String(),
		engine.Bag,
		engine.players,
	}
	return json.Marshal(eng)
}
