package polymers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models/enerr"
	"github.com/google/uuid"
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

// MarshalJSON - метод, который возвращает пустой JSON-объект для объектов SimpleState,
// так как они не имеют внутреннего состояния
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
		return st, enerr.E("Неправильный Action", enerr.InvalidRequest)
	}

	if s.secure[action] {
		if player.Role < common.Judge_Role {
			return st, enerr.E("Недостаточно прав", enerr.Unauthorized)
		}
	}

	handle, ok := s.handlers[action]
	if !ok {
		return st, enerr.E(fmt.Sprintf("Неизвестное действие: %s", action), enerr.NotExistAction)
	}
	st, err := handle(e)
	return st, err
}

func (s SimpleState) PreHook() {
	// PreHook doesnt do anything for SimpleState, becouse its no need for update internals
}

type ObtainState struct {
	SimpleState

	ticker *Ticker
	maxDur time.Duration
	step   int

	eng *PolymersEngine
}

func (eng *PolymersEngine) NewObtainState(timer time.Duration) (state *ObtainState) {
	state = &ObtainState{
		eng:         eng,
		maxDur:      timer,
		step:        0,
		ticker:      newTicker(),
		SimpleState: NewState(),
	}
	state.Add("GetElement", GetElement(eng), true)
	state.Add("RaiseHand", RaiseHand(eng), false)
	return
}

func (s *ObtainState) Handlers() map[string]HandlerFunc {
	return s.handlers
}

func (s *ObtainState) MarshalJSON() ([]byte, error) {
	st := struct {
		Timer int
	}{
		int(s.ticker.Remains().Seconds()),
	}
	return json.Marshal(st)
}
func (s *ObtainState) PreHook() {
	s.handlers["GetElement"](models.Action{})
	s.ticker.Reset(s.incrementalTime())
}

// Update current state - obtain new element from bag.
func (s *ObtainState) Update() (stateInt, error) {
	if !s.ticker.Ticked() {
		return NO_TRANSITION, nil
	}
	s.eng.log.Debug("Updating ObtainState")
	st, err := s.handlers["GetElement"](models.Action{})
	if err != nil {
		return st, err
	}
	if st == OBTAIN || st == NO_TRANSITION {
		s.step += 1
		s.ticker.Reset(s.incrementalTime())
		s.eng.broadcast(common.Message{
			Type: common.ENGINE_ACTION,
			Ok:   true,
			Body: map[string]any{"Action": "NewTimer", "Value": (s.ticker.currentDuration).Seconds()},
		})
	}

	return st, nil
}
func (s *ObtainState) incrementalTime() time.Duration {
	if s.step <= 10 {
		// Для первых 10 ходов вычисляем значение как логарифм по основанию 2
		// Добавляем 1 чтобы избежать логарифма от 0
		dur := time.Duration(math.Log2(float64(s.step)+1)+1) * time.Second
		// fmt.Println(dur.Seconds())
		if dur < time.Duration(s.maxDur) {
			if dur < time.Second {
				dur = time.Second
			}
			return (dur).Round(time.Second)
		}
	}
	// Для всех остальных ходов возвращаем заданное время
	return time.Duration(s.maxDur)

}

type TradeState struct {
	SimpleState
	ticker        *Ticker
	timeForTrade  time.Duration
	eng           *PolymersEngine
	StockExchange *StockExchange
}

func (eng *PolymersEngine) NewTradeState(timer time.Duration) (state *TradeState) {
	state = &TradeState{
		ticker:       newTicker(),
		eng:          eng,
		timeForTrade: timer,
		SimpleState:  NewState(),
		StockExchange: &StockExchange{
			StockList: make([]*Stock, 0, 10),
			Requests:  make(map[string]*StockRequest),
		},
	}
	state.Add("TradeOffer", state.addTradeOffer(), false)
	state.Add("Continue", func(a models.Action) (stateInt, error) { return OBTAIN, nil }, true)

	return

}

func (s *TradeState) PreHook() {
	s.ticker.Reset(s.timeForTrade)
}

// Update current state - switch to Obtain
func (s *TradeState) Update() (stateInt, error) {
	if !s.ticker.Ticked() {
		return NO_TRANSITION, nil
	}
	s.eng.log.Debug("Updating TradeState")
	return OBTAIN, nil
}

type StockExchange struct {
	StockList []*Stock
	Requests  map[string]*StockRequest
}

type StockRequest struct {
	ID     string
	Player string
	Accept bool
}

type Stock struct {
	ID        string
	Owner     *Participant
	Element   string
	ToElement string
}

func (s *StockExchange) AddStock(stck *Stock) {
	s.StockList = append(s.StockList, stck)
}
func (s *StockExchange) AddRequest(StockID string, req *StockRequest) {
	s.Requests[StockID] = req
}

func (s *StockExchange) StockByUser(user string) (*Stock, error) {
	for _, stock := range s.StockList {
		if stock.Owner.Name == user {
			return stock, nil
		}
	}
	return nil, enerr.E("Предложение не найдено", enerr.InvalidRequest)
}

// func (s *StockExchange) SetAck(StockId string, RequestID string){

// }
func (stk *Stock) MarshalJSON() ([]byte, error) {
	st := struct {
		ID        string
		Owner     string
		Element   string
		ToElement string
	}{
		stk.ID, stk.Owner.Name, stk.Element, stk.ToElement,
	}
	return json.Marshal(st)
}

func (s *TradeState) addTradeOffer() HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Element   string
		ToElement string `json:"toElement"`
	}
	return func(e models.Action) (stateInt, error) {

		data, owner, err := dataFromAction[Data](e, s.eng)
		if err != nil {
			return NO_TRANSITION, err
		}

		if owner.Bag[data.Element] < 1 {
			return NO_TRANSITION, enerr.E("У вас нет такого элемента", enerr.GameLogic)
		}

		s.StockExchange.AddStock(&Stock{
			ID:        uuid.NewString(),
			Owner:     owner,
			Element:   data.Element,
			ToElement: data.ToElement,
		})
		s.eng.broadcast(common.Message{
			Type: common.ENGINE_INFO,
			Ok:   true,
			Body: s.eng.PreHook(),
		})
		return NO_TRANSITION, nil
	}
}
func (s *TradeState) addTradeRequest() HandlerFunc {
	type Data struct {
		Type    string
		Action  string
		StockID string
		Accept  bool
	}
	return func(e models.Action) (stateInt, error) {
		data, _, err := dataFromAction[Data](e, s.eng)
		if err != nil {
			return NO_TRANSITION, err
		}
		s.StockExchange.AddRequest(data.StockID, &StockRequest{ID: uuid.NewString(), Player: e.Player, Accept: data.Accept})

		return NO_TRANSITION, nil
	}
}

func (s *TradeState) addTradeAck() HandlerFunc {
	type Data struct {
		Type     string
		Action   string
		TargetID string
	}
	return func(e models.Action) (stateInt, error) {
		_, player, err := dataFromAction[Data](e, s.eng)
		if err != nil {
			return NO_TRANSITION, err
		}
		s.StockExchange.StockByUser(player.Name)
		// s.StockExchange.Requests[]
		// TODO:
		return NO_TRANSITION, nil
	}
}
func (s *TradeState) MarshalJSON() ([]byte, error) {
	st := struct {
		StockExchange StockExchange
		Timer         int
	}{*s.StockExchange, int(s.ticker.Remains())}
	return json.Marshal(st)
}
