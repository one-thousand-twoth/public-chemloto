package polymers

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
)

type HandlerFunc func(models.Action) stateInt

func (engine *PolymersEngine) RaiseHand() HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Field     string
		Name      string
		Structure map[string]int
	}
	return func(e models.Action) stateInt {
		engine.log.Info("Handle RaiseHand", slog.String("player", e.Player), slog.Any("players", engine.players()))
		var data Data
		if err := mapstructure.Decode(e.Envelope, &data); err != nil {
			engine.log.Error("Failed to decode Check data", sl.Err(err))
		}

		player, _ := engine.getPlayer(e.Player)
		var less bool
		for k := range data.Structure {
			less = player.Bag[k] < data.Structure[k]
			if less {
				engine.log.Error("To many elements, than player has")
				engine.unicast(e.Player, common.Message{
					Type:   common.ENGINE_ACTION,
					Ok:     false,
					Errors: []string{"Указаны элементы которых нет в количестве"},
					Body: map[string]any{
						"Action":  "RaiseHand",
						"Players": engine.players(),
					},
				})
				return NO_TRANSITION
			}
		}

		var eq bool
		for _, entry := range engine.Checks.Fields[data.Field][data.Name] {
			eq = reflect.DeepEqual(removeZeroValues(entry), data.Structure)
			if eq {
				break
			}
		}
		if !eq {
			player.Score -= 1
			engine.unicast(player.Name, common.Message{
				Type:   common.UNDEFINED,
				Ok:     false,
				Errors: []string{"Неправильный состав элементов", "Вычитается один балл"},
			})
			engine.broadcast(common.Message{
				Type: common.ENGINE_INFO,
				Ok:   true,
				Body: engine.PreHook(),
			})
			return NO_TRANSITION
		}
		player.RaisedHand = true
		engine.RaisedHands = append(engine.RaisedHands, Hand{Player: player, Field: data.Field, Name: data.Name, Structure: data.Structure})

		engine.broadcast(common.Message{
			Type: common.ENGINE_ACTION,
			Ok:   true,
			Body: map[string]any{
				"Action":  "RaiseHand",
				"Players": engine.players(),
			},
		})
		return HAND
	}
}

func (engine *PolymersEngine) Check() HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Player    string
		Field     string
		Name      string
		Structure map[string]int
	}
	return func(e models.Action) stateInt {
		engine.log.Info("Handle Check action",
			slog.String("player", e.Player),
			slog.Any("players", engine.players()),
			slog.String("data", fmt.Sprintf("%#v", e.Envelope["Structure"])))
		var data Data
		if err := mapstructure.Decode(e.Envelope, &data); err != nil {
			engine.log.Error("Failed to decode Check data", sl.Err(err))
		}
		player, err := engine.getPlayer(data.Player)
		if err != nil {
			engine.log.Error("cannot find player")
			return NO_TRANSITION
		}
		var less bool
		for k := range data.Structure {
			less = player.Bag[k] < data.Structure[k]
			if less {
				engine.log.Error("To many elements, than player has")
				return NO_TRANSITION
			}
		}
		var eq bool
		for _, entry := range engine.Checks.Fields[data.Field][data.Name] {
			eq = reflect.DeepEqual(removeZeroValues(entry), data.Structure)
			if eq {
				break
			}
		}
		if eq {
			engine.log.Info("Succesful check Polymers", slog.String("Player", e.Player), slog.Any("Data", data))
			player.RaisedHand = false
			for k := range player.Bag {
				player.Bag[k] -= data.Structure[k]
			}

		} else {
			engine.log.Error("Failed to check Polymers",
				slog.Any("data", data),
				slog.Any("example", engine.Checks.Fields[data.Field][data.Name]),
			)
			player.RaisedHand = false
			player.Score -= 1
			for i := 0; i < len(engine.RaisedHands); i++ {
				if engine.RaisedHands[i].Player.Name == e.Player {
					engine.RaisedHands = append(engine.RaisedHands[:i], engine.RaisedHands[i+1:]...)
					break
				}
			}
		}
		if len(engine.unchecked()) == 0 {
			engine.log.Debug("all players checked")
			for _, hand := range engine.RaisedHands {
				hand.Player.Score += engine.Fields[hand.Field].getScore()
			}
			engine.RaisedHands = engine.RaisedHands[:0]
			return OBTAIN
		}
		return HAND
	}
}
func (engine *PolymersEngine) GetElement() HandlerFunc {
	return func(e models.Action) stateInt {
		elem, err := engine.Bag.getRandomElement()
		if err != nil {
			if errors.Is(err, ErrEmptyBag) {
				engine.log.Info("Empty bag!")
				elem = "Empty bag!"
				engine.Bag.PushedValues = append(engine.Bag.PushedValues, elem)
			} else {
				engine.log.Error("Error Get Element", sl.Err(err))
				return NO_TRANSITION
			}
		}

		for _, player := range engine.players() {
			player.Bag[elem] += 1
		}
		engine.log.Debug("Got element", slog.String("elem", elem))
		engine.broadcast(common.Message{Type: common.ENGINE_ACTION, Ok: true, Body: map[string]any{
			"Action":       "GetElement",
			"Element":      elem,
			"LastElements": engine.Bag.LastElements(),
		}})
		if elem == "Empty bag!" {
			return COMPLETED
		}
		return NO_TRANSITION
	}

}
