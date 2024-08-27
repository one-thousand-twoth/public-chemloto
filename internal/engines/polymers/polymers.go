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

func New(log *slog.Logger, cfg PolymersEngineConfig) *PolymersEngine {
	src := rand.NewSource(time.Now().UnixNano())

	eng := &PolymersEngine{
		log:    log.With(slog.String("engine", "PolymersEngine")),
		Checks: cfg.Checks,
		Bag:    NewGameBag(cfg.Elements),
		// Fields will be filled latly on [[Start]]
		Fields: map[string]*Field{
			"Альфа": {Score: 0},
			"Бета":  {Score: 0},
			"Гамма": {Score: 0}},
		RaisedHands: make([]Hand, 0),
		ActionChan:  make(chan models.Action),
		// handlers:    map[string]HandlerFunc{},
		// ticker will be filled latly on [[Start]]
		ticker:     &time.Ticker{},
		unicast:    cfg.Unicast,
		broadcast:  cfg.Broadcast,
		timerInt:   cfg.TimerInt,
		MaxPlayers: cfg.MaxPlayers,
		rnd:        rand.New(src),
	}
	eng.StateMachine = stateMachine{
		Current: OBTAIN,
		States: map[stateInt]State{
			OBTAIN: eng.NewObtainState(),
			HAND: NewState().
				Add("RaiseHand", RaiseHand(eng), false).
				Add("Check", Check(eng), true),
		},
	}

	return eng
}

// Use to Configure all Engine params.
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
	Fields       map[string]*Field
	RaisedHands  []Hand
	// Канал для обработки действий игроков
	ActionChan chan models.Action
	// handlers   map[string]HandlerFunc
	Checks   Checks
	Internal chan time.Time
	// TODO: depr
	ticker *time.Ticker

	Participants []*Participant
	MaxPlayers   int

	unicast   unicastFunction
	broadcast broadcastFunction

	timerInt int
	mu       sync.Mutex
	rnd      *rand.Rand
}
type Hand struct {
	Player    *Participant
	Field     string
	Name      string
	Structure map[string]int
}
type Field struct {
	Score int
}

// getter. Score field will be decriadsed by one
func (f *Field) getScore() int {
	o := f.Score
	f.Score -= 1
	return o
}

type Participant struct {
	models.Player
	RaisedHand bool
	Bag        map[string]int
	Score      int
}

// unicastFunction accepts first argument userID
type unicastFunction func(string, common.Message)
type broadcastFunction func(common.Message)

// Start Game
// All inputs should lock engine mutex.
func (engine *PolymersEngine) Start() {
	if engine.started {
		// TODO:
		return
	}
	engine.started = true
	for _, field := range engine.Fields {
		field.Score = len(engine.players())
	}
	if engine.timerInt > 0 {
		engine.ticker = time.NewTicker(time.Duration(engine.timerInt) * time.Second)
	}
	go func() {
		for {
			select {
			case e, _ := <-engine.ActionChan:
				engine.mu.Lock()
				engine.log.Debug("locked engine")
				func() {
					player, err := engine.getPlayer(e.Player)
					if err != nil {
						engine.log.Error("Unknown player", slog.String("name", e.Player))
						return
					}
					state, err := engine.StateMachine.States[engine.StateMachine.Current].Handle(e, player)
					if err != nil {
						if errors.Is(err, ErrNoAuthorized) {
							engine.log.Error("no authorized")
							engine.unicast(e.Player,
								common.Message{Type: common.UNDEFINED,
									Ok:     false,
									Errors: []string{"Недостаточно прав"}})
						}
						engine.log.Error(
							"error while handling action with state",
							sl.Err(err),
							slog.String("state", engine.StateMachine.Current.String()))

						return
					}
					// Changing state if needed and notificate users.
					if state > NO_TRANSITION {
						engine.log.Info("Changing game state",
							slog.String("new state", state.String()),
							slog.String("old state", engine.StateMachine.Current.String()))
						engine.StateMachine.Current = state
						engine.StateMachine.States[state].PreHook()
						engine.broadcast(common.Message{
							Type: common.ENGINE_INFO,
							Ok:   true,
							Body: engine.PreHook(),
						})
					}
				}()
				engine.mu.Unlock()
				engine.log.Debug("unlocked engine")
			case _ = <-engine.ticker.C:
				engine.mu.Lock()
				func() {
					// Избыточная проверка, потому что предполагаю, что
					// есть маленький шанс, когда тик может прийти после смены state
					if engine.StateMachine.Current != OBTAIN {
						return
					}
					engine.StateMachine.States[OBTAIN].Handlers()["GetElement"](models.Action{})

				}()
				engine.mu.Unlock()
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
		return models.ErrAlreadyStarted
	}
	if player.Role == common.Player_Role {
		if len(engine.players()) >= engine.MaxPlayers {
			return models.ErrMaxPlayers
		}
	}
	engine.Participants = append(engine.Participants, &Participant{Player: player, Bag: make(map[string]int)})
	return nil
}

func (engine *PolymersEngine) RemovePlayer(name string) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	if engine.started {
		return models.ErrAlreadyStarted
	}
	for i := 0; i < len(engine.Participants); i++ {
		if engine.Participants[i].Name == name {
			engine.Participants = append(engine.Participants[:i], engine.Participants[i+1:]...)
			break
		}
	}
	return nil
}

// getPlayer internal. Not conccurent safe
func (engine *PolymersEngine) getPlayer(name string) (*Participant, error) {
	for i := 0; i < len(engine.Participants); i++ {
		if engine.Participants[i].Name == name {
			return engine.Participants[i], nil
		}
	}
	return &Participant{}, errors.New("unknown player")
}

// players() return engine.Participants with Player Role
func (engine *PolymersEngine) players() []*Participant {
	players := make([]*Participant, 0, len(engine.Participants))
	for i := 0; i < len(engine.Participants); i++ {
		if engine.Participants[i].Role == common.Player_Role {
			players = append(players, engine.Participants[i])
		}
	}
	return players
}

// Возвращает список игроков которых еще не проверяли
func (engine *PolymersEngine) unchecked() []*Participant {
	players := make([]*Participant, 0, len(engine.Participants))
	for i := 0; i < len(engine.Participants); i++ {
		if engine.Participants[i].Role == common.Player_Role && engine.Participants[i].RaisedHand {
			players = append(players, engine.Participants[i])
		}
	}
	return players
}

// Control what data will be send on [ENGINE_INFO] Message and etc.
func (engine *PolymersEngine) MarshalJSON() ([]byte, error) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	players := make([]Participant, len(engine.Participants))
	fields := make(map[string]Field, len(engine.Fields))

	// Проходим по каждому элементу и копируем его в новый срез
	for i, p := range engine.Participants {
		players[i] = *p
	}
	for k, v := range engine.Fields {
		fields[k] = *v
	}
	eng := struct {
		Started     bool
		State       string
		Bag         GameBag
		Players     []Participant
		RaisedHands []Hand
		Fields      map[string]Field
	}{
		engine.started,
		engine.StateMachine.Current.String(),
		engine.Bag,
		players,
		engine.RaisedHands,
		fields,
	}
	return json.Marshal(eng)
}
