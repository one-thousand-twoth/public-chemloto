package polymers

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
)

type HandlerFunc func(models.Action) (stateInt, error)

func RaiseHand(engine *PolymersEngine) HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Field     string
		Name      string
		Structure map[string]int
	}
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "polymers/RaiseHand"
		data, err := dataFromAction[Data](e, engine)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}
		player, err := engine.getPlayer(e.Player)
		if err != nil {
			return NO_TRANSITION, enerr.E(op, err)
		}

		engine.log.Debug("Handle RaiseHand",
			slog.String("player", e.Player),
			slog.Any("players", engine.players()),
			enerr.OpAttr(op))

		var less bool
		for k := range data.Structure {
			less = player.Bag[k] < data.Structure[k]
			if less {
				return NO_TRANSITION, enerr.E(op, "Указаны элементы которых нет в таком количестве", enerr.GameLogic)
			}
		}

		var eq bool
		for _, entry := range engine.checks.Fields[data.Field][data.Name] {
			eq = reflect.DeepEqual(removeZeroValues(entry), data.Structure)
			if eq {
				break
			}
		}
		if !eq {
			player.score(-1)
			engine.broadcast(common.Message{
				Type: common.ENGINE_INFO,
				Ok:   true,
				Body: engine.PreHook(),
			})
			return NO_TRANSITION, enerr.E(op, "Неправильный состав элементов")
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
		Type      string
		Action    string
		Player    string
		Field     string
		Name      string
		Structure map[string]int
	}
	return func(e models.Action) (stateInt, error) {

		const op enerr.Op = "polymers/Check"
		data, err := dataFromAction[Data](e, engine)
		if err != nil {
			return NO_TRANSITION, err
		}
		target, err := engine.getPlayer(data.Player)
		if err != nil {
			return NO_TRANSITION, err
		}

		engine.log.Debug("Handle Check action",
			slog.String("from_player", e.Player),
			slog.String("target", data.Player),
			slog.Any("players", engine.players()),
			slog.String("Field", data.Field),
			slog.String("Name", data.Name),
			slog.Any("Structure", data.Structure),
			enerr.OpAttr(op),
		)
		// TODO: роазобраться, тут кажется напутана логика
		var less bool
		for k := range data.Structure {
			less = target.Bag[k] < data.Structure[k]
			if less {
				engine.log.Error("To many elements, than player has")
				return NO_TRANSITION, enerr.E(op, "Слишком много элементов", enerr.GameLogic)
			}
		}
		var eq bool
		// Смотрим достаточно ли элементов у игрока в мешке
		for _, entry := range engine.checks.Fields[data.Field][data.Name] {
			eq = reflect.DeepEqual(removeZeroValues(entry), data.Structure)
			if eq {
				break
			}
		}
		// TODO: обработать ошибку
		if eq {
			engine.log.Info("Succesful check Polymers",
				slog.String("Player", e.Player),
				slog.Any("Data", data),
				enerr.OpAttr(op))
			target.RaisedHand = false
			for k := range target.Bag {
				target.Bag[k] -= data.Structure[k]
			}

		} else {
			engine.log.Error("Failed to check Polymers",
				slog.Any("data", data),
				slog.Any("example", engine.checks.Fields[data.Field][data.Name]),
				enerr.OpAttr(op),
			)
			target.RaisedHand = false
			target.score(-1)
			for i := 0; i < len(engine.raisedHands); i++ {
				if engine.raisedHands[i].Player.Name == e.Player {
					engine.raisedHands = append(engine.raisedHands[:i], engine.raisedHands[i+1:]...)
					break
				}
			}

		}
		if len(engine.unchecked()) == 0 {
			engine.log.Debug("all players checked", enerr.OpAttr(op))
			for _, hand := range engine.raisedHands {
				hand.Player.score(engine.fields[hand.Field].decrementScore())
			}
			engine.raisedHands = engine.raisedHands[:0]
			return OBTAIN, nil
		}
		return HAND, nil
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
			return COMPLETED, nil
		}
		if elem == "TRADE" {
			return TRADE, nil
		}
		return NO_TRANSITION, nil
	}
}
func (engine *PolymersEngine) Trade() HandlerFunc {
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
		pl1, err := engine.getPlayer(data.Player1)
		if err != nil {
			return NO_TRANSITION, err
		}
		pl2, err := engine.getPlayer(data.Player2)
		if err != nil {
			return NO_TRANSITION, err
		}
		if pl1.Bag[data.Element1] <= 1 {
			// engine.log.Error("found no element", "user", pl1, "element", data.Element1)
			return NO_TRANSITION, enerr.E(fmt.Sprintf("У игрока %s нет такого элемента", pl1.Name))
		}
		if pl2.Bag[data.Element2] <= 1 {
			// engine.log.Error("found no element", "user", pl1, "element", data.Element1)
			return NO_TRANSITION, enerr.E(fmt.Sprintf("У игрока %s нет такого элемента", pl2.Name))
		}
		pl1.Bag[data.Element1] -= 1
		pl2.Bag[data.Element2] -= 1
		pl1.Bag[data.Element2] += 1
		pl2.Bag[data.Element1] += 1
		return TRADE, nil
	}
}
