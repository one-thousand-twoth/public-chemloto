package polymers

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
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
	Elements  map[string]int
	TimerInt  int
	Unicast   unicastFunction
	Broadcast broadcastFunction
}

type PolymersEngine struct {
	log          *slog.Logger
	started      bool
	StateMachine stateMachine
	Bag          GameBag
	// Канал для обработки действий игроков
	ActionChan chan models.Action
	handlers   map[string]HandlerFunc

	players []string

	unicast   unicastFunction
	broadcast broadcastFunction

	timerInt int

	rnd *rand.Rand
}

// unicastFunction accepts first argument userID
type unicastFunction func(string, common.Message)
type broadcastFunction func(common.Message)

// TODO: OnlyOnce запустить только раз.
func (engine *PolymersEngine) Start() {
	engine.started = true
	go func() {
		for e := range engine.ActionChan {
			state, err := engine.StateMachine.States[engine.StateMachine.Current].Input(e)
			if err != nil {
				engine.log.Error("error while handling action with state %s", sl.Err(err), engine.StateMachine.Current)
			}
			// Changing state if needed
			if state > NO_TRANSITION {
				engine.StateMachine.Current = state
			}
		}
	}()
	engine.log.Debug("Broadcast for starting engine")
	engine.broadcast(common.Message{Type: common.HUB_STARTGAME, Ok: true})
}

func (engine *PolymersEngine) Input(e models.Action) {
	engine.ActionChan <- e
}
func (engine *PolymersEngine) PreHook() {

}

func (engine *PolymersEngine) AddPlayer(name string) {
	if engine.started {
		return
	}
	engine.players = append(engine.players, name)
}

func (engine *PolymersEngine) RemovePlayer(name string) {
	if engine.started {
		return
	}
	engine.players = removeByValue(engine.players, name)
}
