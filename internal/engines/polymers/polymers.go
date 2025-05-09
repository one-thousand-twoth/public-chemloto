package polymers

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
)

// New returns PolymersEngine with cfg parameters
func New(log *slog.Logger, cfg PolymersEngineConfig) *PolymersEngine {
	src := rand.NewSource(time.Now().UnixNano())

	eng := &PolymersEngine{
		log:    log.With(slog.String("engine", "PolymersEngine")),
		checks: cfg.Checks,
		bag:    NewGameBag(cfg.Elements),
		// Fields will be filled lately on [[Start]]
		fields: map[string]*Field{
			"Альфа": {Score: 0},
			"Бета":  {Score: 0},
			"Гамма": {Score: 0}},
		raisedHands: make([]Hand, 0),
		actionChan:  make(chan models.Action),
		unicast:     cfg.Unicast,
		broadcast:   cfg.Broadcast,
		maxPlayers:  cfg.MaxPlayers,
		done:        make(chan struct{}, 1),
		rnd:         rand.New(src),
	}

	var obtainState State
	if cfg.TimerInt > 0 {
		obtainState = eng.NewObtainState(time.Second*time.Duration(cfg.TimerInt), cfg.IsAutoCheck)
	}
	if cfg.TimerInt == 0 {
		obtainState = NewState().
			Add("GetElement", GetElement(eng), true).
			Add("RaiseHand", RaiseHand(eng, cfg.IsAutoCheck), false).
			Add("AddScore", AddScore(eng), true)
	}
	// Конфигурация FSM и его обработчиков событий.
	eng.stateMachine = stateMachine{
		Current: OBTAIN,
		States: map[stateInt]State{
			OBTAIN: obtainState,
			HAND: NewState().
				Add("RaiseHand", RaiseHand(eng, cfg.IsAutoCheck), false).
				Add("Check", Check(eng), true).
				Add("AddScore", AddScore(eng), true),
			// TRADE: NewState().
			// 	Add("Trade", eng.Trade(), true).
			// 	Add("Continue", func(a models.Action) stateInt { return OBTAIN }, true),
			TRADE:     eng.NewTradeState(time.Hour * 30),
			COMPLETED: NewState(),
		},
	}

	return eng
}

// Use to Configure Engine params.
type PolymersEngineConfig struct {
	Elements map[string]int
	Checks   Checks
	// Should be >= 0
	TimerInt   int
	Unicast    models.UnicastFunction
	Broadcast  models.BroadcastFunction
	MaxPlayers int
	// Проверять ли структуру игрока при поднятии руки
	IsAutoCheck bool
}

// For test cases
func (c PolymersEngineConfig) ToMap() map[string]any {
	cjson, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	cmap := make(map[string]any)
	err = json.Unmarshal(cjson, &cmap)
	if err != nil {
		panic(err)
	}
	return cmap
}

type PolymersEngine struct {
	log    *slog.Logger
	status models.EngineStatus
	// FSM для управления игровыми состояниями
	stateMachine stateMachine
	// Игровой мешок
	bag GameBag
	// Игровые поля представленные в игре
	fields map[string]*Field
	// Список поднятых рук игроков
	raisedHands []Hand
	// Канал для обработки действий игроков
	actionChan chan models.Action
	// Для проверки правильности собранных структур
	checks Checks
	// NOTE: пока не реализованная идея
	// internal chan time.Time

	// Список участников
	participants []*Player
	// Максимум игроков
	maxPlayers int

	// Функция для отправки сообщения одному адресату
	unicast models.UnicastFunction
	// Функция для отправки сообщения всем в комнате
	broadcast models.BroadcastFunction

	done chan struct{}
	mu   sync.Mutex
	rnd  *rand.Rand
}

// Start Game.
//
// All inputs must lock engine mutex.
func (engine *PolymersEngine) Start() {
	if engine.status != models.STATUS_WAITING {
		// TODO: вернуть ошибку
		return
	}
	engine.status = models.STATUS_STARTED
	for _, field := range engine.fields {
		field.Score = len(engine.players())
	}
	engine.stateMachine.States[engine.stateMachine.Current].PreHook()
	go func() {
		for {
			select {
			case e := <-engine.actionChan:
				engine.mu.Lock()
				engine.log.Debug("locked engine")
				func() {
					player, err := engine.getParticipant(e.Player)
					if err != nil {
						enerr.ErrorResponse(engine.unicast, e.Player, engine.log, err)
						return
					}
					state, err := engine.stateMachine.States[engine.stateMachine.Current].Handle(e, player)
					if err != nil {
						enerr.ErrorResponse(engine.unicast, e.Player, engine.log, err)
						// return
					}
					if state == UPDATE_CURRENT {
						engine.broadcast(common.Message{
							Type: common.ENGINE_INFO,
							Ok:   true,
							Body: engine.PreHook(),
						})
					}
					// Changing state if needed and notificate users.
					if state > NO_TRANSITION {
						engine.log.Info("Changing game state",
							slog.String("new state", state.String()),
							slog.String("old state", engine.stateMachine.Current.String()))
						engine.stateMachine.Current = state
						engine.stateMachine.States[state].PreHook()
						engine.broadcast(common.Message{
							Type: common.ENGINE_INFO,
							Ok:   true,
							Body: engine.PreHook(),
						})
					}
				}()
				engine.mu.Unlock()
				engine.log.Debug("unlocked engine")
			case <-engine.done:
				engine.mu.Lock()
				engine.status = models.STATUS_COMPLETED
				engine.broadcast(common.Message{
					Type: common.ENGINE_INFO,
					Ok:   true,
					Body: engine.PreHook(),
				})
				engine.mu.Unlock()
				return
			default:
				engine.mu.Lock()
				state, err := engine.stateMachine.States[engine.stateMachine.Current].Update()
				if err != nil {
					engine.log.Error(
						"error while updating state",
						sl.Err(err),
						slog.String("state", engine.stateMachine.Current.String()))
					break
				}
				if state > NO_TRANSITION {
					engine.log.Info("Changing game state by Update",
						slog.String("new state", state.String()),
						slog.String("old state", engine.stateMachine.Current.String()))
					engine.stateMachine.Current = state
					engine.stateMachine.States[state].PreHook()
					engine.broadcast(common.Message{
						Type: common.ENGINE_INFO,
						Ok:   true,
						Body: engine.PreHook(),
					})
				}
				engine.mu.Unlock()
			}
		}
	}()
	engine.log.Debug("Broadcast for starting engine")
	engine.broadcast(common.Message{Type: common.HUB_STARTGAME, Ok: true})
	engine.broadcast(common.Message{Type: common.ENGINE_INFO, Ok: true, Body: engine.PreHook()})
}
func (engine *PolymersEngine) Exit() {
	if engine.status != models.STATUS_STARTED {
		return
	}
	engine.done <- struct{}{}
}

// Input will send action to engine for processing
func (engine *PolymersEngine) Input(e models.Action) {
	if engine.status != models.STATUS_STARTED {
		engine.unicast(e.Player, common.Message{
			Type:   common.UNDEFINED,
			Ok:     false,
			Errors: []string{"Игра ещё не начата"},
		})
		return
	}
	engine.actionChan <- e
}

// PreHook returns map with pointer to the engine:
//
//	map[string]any{"engine": engine}
//
// Should only be used for MarshalJSON as it is safe for concurent use
func (engine *PolymersEngine) PreHook() map[string]any {
	return map[string]any{"engine": engine}
}

func (engine *PolymersEngine) GetResults() [][]string {
	results := [][]string{{"Игрок", "Очки"}}
	for _, v := range engine.players() {
		results = append(results, []string{v.Name, strconv.Itoa(v.Score)})
	}
	engine.log.Debug("New Results", slog.Any("res", results))
	return results
}

func (engine *PolymersEngine) AddParticipant(player models.Participant) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	if engine.status == models.STATUS_STARTED {
		return enerr.E("Игрок не может быть добавлен так как игра уже начата", enerr.AlreadyStarted)
	}
	if player.Role == common.Player_Role {
		if len(engine.players()) >= engine.maxPlayers {
			return enerr.E("Недостаточно прав", enerr.MaxPlayers)
		}
	}
	engine.participants = append(engine.participants,
		&Player{
			Participant:     player,
			Bag:             map[string]int{"H": 50, "C": 50, "C6H5": 50, "CH3": 40},
			CompletedFields: make([]string, 0, engine.maxPlayers),
		},
	)
	return nil
}

func (engine *PolymersEngine) RemoveParticipant(name string) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	if engine.status == models.STATUS_STARTED {
		return enerr.E("Игрок не может быть удалён так как игра уже начата", enerr.AlreadyStarted)
	}
	for i := 0; i < len(engine.participants); i++ {
		if engine.participants[i].Name == name {
			engine.participants = append(engine.participants[:i], engine.participants[i+1:]...)
			break
		}
	}
	return nil
}

// getParticipant internal. Not concurrent safe.
//
// Can return:
//
//	enerr.Unidentified
func (engine *PolymersEngine) getParticipant(name string) (*Player, error) {
	const op enerr.Op = "polymers/PolymersEngine.getPlayer"
	for i := 0; i < len(engine.participants); i++ {
		if engine.participants[i].Name == name {
			return engine.participants[i], nil
		}
	}
	return &Player{}, enerr.E("Игрок с таким именем не найден", enerr.Unidentified, op)
}

// players() return engine.Participants with Player Role
func (engine *PolymersEngine) players() []*Player {
	players := make([]*Player, 0, engine.maxPlayers)
	for i := 0; i < len(engine.participants); i++ {
		if engine.participants[i].Role == common.Player_Role {
			players = append(players, engine.participants[i])
		}
	}
	return players
}

// Возвращает список игроков которых еще не проверяли после RaiseHand
func uncheckedPlayers(players []*Player) []*Player {
	unchecked := make([]*Player, 0, len(players))
	for i := 0; i < len(players); i++ {
		if players[i].Role == common.Player_Role && players[i].RaisedHand {
			unchecked = append(unchecked, players[i])
		}
	}
	return unchecked
}

// Control what data will be send on [ENGINE_INFO] Message and etc.
func (engine *PolymersEngine) MarshalJSON() ([]byte, error) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	players := make([]Player, len(engine.participants))
	fields := make(map[string]Field, len(engine.fields))

	// Проходим по каждому элементу и копируем его в новый срез
	for i, p := range engine.participants {
		players[i] = *p
	}
	for k, v := range engine.fields {
		fields[k] = *v
	}
	eng := struct {
		Status      string
		State       string
		Bag         GameBag
		Players     []Player
		RaisedHands []Hand
		Fields      map[string]Field
		StateStruct State
	}{
		engine.status.String(),
		engine.stateMachine.Current.String(),
		engine.bag,
		players,
		engine.raisedHands,
		fields,
		engine.stateMachine.States[engine.stateMachine.Current],
	}
	return json.Marshal(eng)
}

// dataFromAction is a convenience function for extracting data from action,
// including Participant information
//
// it can return enerr.InvalidRequest, enerr.Unidentified
func dataFromAction[T any](e models.Action) (*T, error) {
	const op enerr.Op = "polymers/dataFromAction"
	var data T
	if err := mapstructure.Decode(e.Envelope, &data); err != nil {
		return nil, enerr.E("Неправильный запрос", enerr.InvalidRequest, err, op)
	}
	return &data, nil
}
