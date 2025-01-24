package polymers

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
)

type HandlerFunc func(models.Action) (stateInt, error)

func RaiseHand(engine *PolymersEngine, isAuto bool) HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Field     string
		Name      string
		Structure map[string]int
	}
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "polymers/RaiseHand"
		data, err := dataFromAction[Data](e)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		player, err := engine.getParticipant(e.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}

		_, ok := lo.Find(engine.raisedHands, func(v Hand) bool {
			return v.Player.Name == player.Name && v.Field == data.Field
		})
		if ok {
			return NO_TRANSITION, enerr.E(op, "Вы уже подняли руку", enerr.GameLogic)
		}

		engine.log.Debug("Handle RaiseHand",
			slog.String("player", e.Player),
			slog.Any("players", engine.players()),
			enerr.OpAttr(op))

		if invalid, err := player.checkIfHasElements(data.Structure); err != nil {
			return NO_TRANSITION, enerr.E(op, "У вас недостаточно элементов для этой структуры",
				enerr.GameLogic, enerr.Parameter(fmt.Sprintf("%+v", invalid)))
		}
		if isAuto {
			eq := checkFields(engine.checks, data.Field, data.Name, data.Structure)
			if !eq {
				player.setScore(-1)
				return UPDATE_CURRENT, enerr.E(op, "Неправильный состав элементов")
			}
		}
		player.RaisedHand = true
		engine.raisedHands = append(engine.raisedHands,
			Hand{Player: player, Field: data.Field, Name: data.Name, Structure: data.Structure})

		engine.broadcast(common.Message{
			Type: common.ENGINE_ACTION,
			Ok:   true,
			Body: map[string]any{
				"Action":  "RaiseHand",
				"Players": engine.players(),
			},
		})
		return HAND, nil
	}
}

func Check(engine *PolymersEngine) HandlerFunc {
	type Data struct {
		Type   string
		Action string
		Accept bool
		Player string
	}
	return func(e models.Action) (stateInt, error) {

		const op enerr.Op = "polymers/Check"
		data, err := dataFromAction[Data](e)
		if err != nil {
			return NO_TRANSITION, err
		}
		target, err := engine.getParticipant(data.Player)
		if err != nil {
			return NO_TRANSITION, err
		}

		engine.log.Debug("Handle Check action",
			slog.String("from_player", e.Player),
			slog.String("target", data.Player),
			slog.Any("players-not-checked", uncheckedPlayers(engine.players())),
			slog.Any("players-not-checked", engine.players()),

			enerr.OpAttr(op),
		)
		hand, index, ok := lo.FindIndexOf(engine.raisedHands, func(v Hand) bool { return v.Player.Name == data.Player })
		if !ok {
			return NO_TRANSITION, enerr.E(op, "Не найдено", enerr.NotExist)
		}
		if hand.Checked {
			return NO_TRANSITION, enerr.E(op, "Вы уже проверили этого игрока", enerr.GameLogic)
		}
		if !data.Accept {
			target.setScore(-1)
			// NOTE: Не делаю zeroing elements, так как после перехода на новые версии компилятора
			// DeleteFunc сам будет занулять их. Пока влиянием на производительность можно принебречь.
			engine.raisedHands = slices.DeleteFunc(engine.raisedHands, func(v Hand) bool {
				return v.Player.Name == hand.Player.Name
			})
		} else {
			engine.raisedHands[index].Checked = true
			target.CompletedFields = append(target.CompletedFields, hand.Field)
			for k := range target.Bag {
				target.Bag[k] -= hand.Structure[k]
			}
		}
		target.RaisedHand = false

		engine.log.Info("Succesful check Polymers",
			slog.String("Player", e.Player),
			slog.Any("Data", data),
			enerr.OpAttr(op))

		if len(uncheckedPlayers(engine.players())) == 0 {
			engine.log.Debug("all players checked", enerr.OpAttr(op))
			for _, hand := range engine.raisedHands {
				hand.Player.setScore(engine.fields[hand.Field].decrementScore())
			}
			// clear slice, without changing capacity
			engine.raisedHands = engine.raisedHands[:0]
			return OBTAIN, nil
		}
		return UPDATE_CURRENT, nil
	}
}
func GetElement(engine *PolymersEngine) HandlerFunc {
	return func(_ models.Action) (stateInt, error) {
		const op enerr.Op = "polymers/PolymersEngine.GetElement"
		empty := "Empty bag!"
		elem, err := engine.bag.getRandomElement()
		if err != nil {
			// TODO: ErrEmptyBag не должен быть ошибкой, по идее - надо поправить
			if errors.Is(err, ErrEmptyBag) {
				engine.log.Info("Элементы в мешке закончились", enerr.OpAttr(op))
				elem = empty
				engine.bag.PushedValues = append(engine.bag.PushedValues, elem)
			} else {
				engine.log.Error("Error Get Element", sl.Err(err))
				return NO_TRANSITION, err
			}
		}

		for _, player := range engine.players() {
			player.Bag[elem] += 1
		}
		engine.log.Debug("Got new element", slog.String("elem", elem), enerr.OpAttr(op))
		engine.broadcast(common.Message{Type: common.ENGINE_ACTION, Ok: true, Body: map[string]any{
			"Action":       "GetElement",
			"Element":      elem,
			"LastElements": engine.bag.LastElements(),
		}})
		if elem == empty {
			engine.done <- struct{}{}
			return COMPLETED, nil
		}
		if elem == "TRADE" {
			return TRADE, nil
		}
		return NO_TRANSITION, nil
	}
}
func (engine *PolymersEngine) TradeHandler() HandlerFunc {
	type Data struct {
		Type     string
		Action   string
		Player1  string
		Element1 string
		Player2  string
		Element2 string
	}
	return func(e models.Action) (stateInt, error) {
		var data Data
		if err := mapstructure.Decode(e.Envelope, &data); err != nil {
			return NO_TRANSITION, err
		}
		pl1, err := engine.getParticipant(data.Player1)
		if err != nil {
			return NO_TRANSITION, err
		}
		pl2, err := engine.getParticipant(data.Player2)
		if err != nil {
			return NO_TRANSITION, err
		}

		err = engine.exchange(pl1, data.Element1, data.Element2, pl2)
		if err != nil {
			return NO_TRANSITION, err
		}
		return TRADE, nil
	}
}

// exchange will exchange elements between players
//
// element1 is player1 element, element2 is player2 element
func (*PolymersEngine) exchange(
	player1 *Player,
	element1 string,
	element2 string,
	player2 *Player) error {
	const op enerr.Op = "polymers/PolymersEngine.exchange"
	if player1.Bag[element1] < 1 {
		return enerr.E(op, fmt.Sprintf("У игрока %s нет такого элемента", player1.Name))
	}
	if player2.Bag[element2] < 1 {
		return enerr.E(op, fmt.Sprintf("У игрока %s нет такого элемента", player2.Name))
	}
	player1.Bag[element1] -= 1
	player2.Bag[element2] -= 1
	player1.Bag[element2] += 1
	player2.Bag[element1] += 1
	return nil
}
