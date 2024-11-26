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
	// for internal use. Signal engine for notify users about current state update
	UPDATE_CURRENT stateInt = -1

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
		},
	}
	state.Add("TradeOffer", state.addTradeOffer(), false)
	state.Add("RemoveTradeOffer", state.removeTradeOffer(), false)

	state.Add("TradeRequest", state.addTradeRequest(), false)
	state.Add("TradeAck", state.addTradeAck(), false)
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
	StockList     []*Stock
	ExchangedList map[string]bool
	// Requests  map[string][]*StockRequest
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
	Requests  map[string]*StockRequest
}

func (stk *Stock) MarshalJSON() ([]byte, error) {
	st := struct {
		ID        string
		Owner     string
		Element   string
		ToElement string
		Requests  map[string]*StockRequest
	}{
		stk.ID, stk.Owner.Name, stk.Element, stk.ToElement, stk.Requests,
	}
	return json.Marshal(st)
}
func (stk *Stock) Request(id string) (*StockRequest, error) {
	const op enerr.Op = "polymers/Stock.Request"
	for _, request := range stk.Requests {
		if request.ID == id {
			return request, nil
		}
	}
	return nil, enerr.E(op, "Предложение не найдено", enerr.InvalidRequest)
}

func (s *StockExchange) AddStock(id string, stck *Stock) {
	if stck.Requests == nil {
		stck.Requests = make(map[string]*StockRequest, 0)
	}
	s.StockList = append(s.StockList, stck)
}

func (s *StockExchange) StockByID(id string) (*Stock, error) {
	const op enerr.Op = "polymers/StockExchange.StockByID"
	for _, stock := range s.StockList {
		if stock.ID == id {
			return stock, nil
		}
	}
	return nil, enerr.E(op, "Предложение не найдено", enerr.InvalidRequest)
}

func (s *StockExchange) StockByUser(user string) (*Stock, error) {
	const op enerr.Op = "polymers/StockExchange.StockByUser"
	for _, stock := range s.StockList {
		if stock.Owner.Name == user {
			return stock, nil
		}
	}
	return nil, enerr.E(op, "Предложение не найдено", enerr.InvalidRequest)
}

func (s *StockExchange) RemoveStockByUser(user string) error {
	const op enerr.Op = "polymers/StockExchange.RemoveStockByUser"
	for i, stock := range s.StockList {
		if stock.Owner.Name == user {
			s.StockList = append(s.StockList[:i], s.StockList[i+1:]...)
			return nil
		}
	}
	return enerr.E(op, "Предложение не найдено", enerr.InvalidRequest)
}

func (s *StockExchange) SetRequest(stock string, req *StockRequest) error {
	const op enerr.Op = "polymers/StockExchange.SetRequest"
	for _, stck := range s.StockList {
		if stck.ID == stock {
			stck.Requests[req.Player] = req
			return nil
		}
	}
	return enerr.E(op, "Предложение не найдено", enerr.NotExist)
}
func (s *StockExchange) SetUserExchanged(user string) error {
	const op enerr.Op = "polymers/StockExchange.SetUserExchanged"
	_, ok := s.ExchangedList[user]
	if ok {
		return enerr.E(op, "Пользователь уже обменялся ранее", enerr.GameLogic)
	}
	s.ExchangedList[user] = true
	return nil
}
func (s *StockExchange) isUserExchanged(user string) bool {
	const op enerr.Op = "polymers/StockExchange.isUserExchanged"
	_, ok := s.ExchangedList[user]
	return ok
}

func (s *StockExchange) SetAck(stockId string, RequestID string) error {
	const op enerr.Op = "polymers/StockExchange.SetAck"
	for _, stck := range s.StockList {
		if stck.ID == stockId {
			// if _, ok := stck.Requests[req.Player]; ok {
			// 	return enerr.E(op, "Пользователь уже дал ответ", enerr.InvalidRequest)
			// }
			req, ok := stck.Requests[RequestID]
			if ok {
				if !req.Accept {
					return enerr.E(op, "Предложение недействительно")
				}

			}
		}
	}
	return enerr.E(op, "Предложение не найдено", enerr.NotExist)
}

func (s *TradeState) addTradeOffer() HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Element   string
		ToElement string `json:"toElement"`
	}
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "polymers/TradeState.addTradeOffer"
		data, err := dataFromAction[Data](e)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		owner, err := s.eng.getPlayer(e.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		if s.StockExchange.isUserExchanged(owner.Name) {
			return NO_TRANSITION, enerr.E(op, "Вы уже обменялись на этом ходу", enerr.GameLogic)
		}
		if owner.Bag[data.Element] < 1 {
			return NO_TRANSITION, enerr.E(op, "У вас нет такого элемента", enerr.GameLogic)
		}
		uuid := uuid.NewString()
		s.StockExchange.AddStock(uuid, &Stock{
			ID:        uuid,
			Owner:     owner,
			Element:   data.Element,
			ToElement: data.ToElement,
			Requests:  make(map[string]*StockRequest),
		})
		return UPDATE_CURRENT, nil
	}
}
func (s *TradeState) removeTradeOffer() HandlerFunc {
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "polymers/TradeState.removeTradeOffer"
		player, err := s.eng.getPlayer(e.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		err = s.StockExchange.RemoveStockByUser(player.Name)
		if err != nil {
			return NO_TRANSITION, err
		}

		return UPDATE_CURRENT, nil
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
		const op enerr.Op = "polymers/TradeState.addTradeRequest"
		data, err := dataFromAction[Data](e)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		player, err := s.eng.getPlayer(e.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		if s.StockExchange.isUserExchanged(player.Name) {
			return NO_TRANSITION, enerr.E(op, "Вы уже обменялись на этом ходу", enerr.GameLogic)
		}
		err = s.StockExchange.SetRequest(
			data.StockID,
			&StockRequest{ID: uuid.NewString(), Player: player.Name, Accept: data.Accept},
		)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		return UPDATE_CURRENT, nil
	}
}

func (s *TradeState) addTradeAck() HandlerFunc {
	type Data struct {
		Type     string
		Action   string
		TargetID string
	}
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "polymers/TradeState.addTradeAck"
		player, err := s.eng.getPlayer(e.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		if s.StockExchange.isUserExchanged(player.Name) {
			return NO_TRANSITION, enerr.E(op, "Вы уже обменялись на этом ходу", enerr.GameLogic)
		}
		data, err := dataFromAction[Data](e)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		stock, err := s.StockExchange.StockByUser(player.Name)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		request, err := stock.Request(data.TargetID)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}

		if !request.Accept {
			return NO_TRANSITION, enerr.E(op, "Игрок не хочет меняться", enerr.GameLogic)
		}

		playerThatRequest, err := s.eng.getPlayer(request.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		if s.StockExchange.isUserExchanged(player.Name) {
			return NO_TRANSITION, enerr.E(op, "Этот игрок уже обменялся с кем то другим", enerr.GameLogic)
		}

		if err := s.eng.exchange(player, stock.Element, stock.ToElement, playerThatRequest); err != nil {
			// по идее ошибки быть не может, но пусть будет проверка
			return NO_TRANSITION, enerr.E(op, err)
		}
		s.StockExchange.RemoveStockByUser(player.Name)
		return UPDATE_CURRENT, nil
	}
}

func (s *TradeState) MarshalJSON() ([]byte, error) {
	st := struct {
		StockExchange StockExchange
		Timer         int
	}{*s.StockExchange, int(s.ticker.Remains().Seconds())}
	return json.Marshal(st)
}
