package polymers

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
)

type HandlerFunc func(models.Action) stateInt

func (engine *PolymersEngine) RaiseHand() HandlerFunc {
	return func(e models.Action) stateInt {
		engine.log.Info("Handle RaiseHand", slog.String("player", e.Player), slog.Any("players", engine.players))
		for k, v := range engine.players {
			if v.Name == e.Player {
				v.RaisedHand = true
				engine.players[k] = v
			}
		}
		engine.broadcast(common.Message{
			Type: common.ENGINE_ACTION,
			Ok:   true,
			Body: map[string]any{
				"Action":  "RaiseHand",
				"Players": engine.players,
			},
		})
		return HAND
	}
}

func (engine *PolymersEngine) Check() HandlerFunc {
	type Data struct {
		Type      string
		Action    string
		Field     string
		Name      string
		Structure map[string]int
	}
	return func(e models.Action) stateInt {
		engine.log.Info("Handle Check action", slog.String("player", e.Player), slog.Any("players", engine.players))
		var data Data
		if err := mapstructure.Decode(e.Envelope, &data); err != nil {
			engine.log.Error("Failed to decode Check data", sl.Err(err))
		}
		var eq bool
		for _, entry := range engine.Checks.Fields[data.Field][data.Name] {
			eq = reflect.DeepEqual(entry, data.Structure)
			if eq {
				break
			}
		}
		if eq {
			engine.log.Info("Succesful check Polymers", slog.String("Player", e.Player), slog.Any("Data", data))
		} else {
			engine.log.Error("Failed to check Polymers", slog.Any("data", data), slog.Any("example", engine.Checks.Fields[data.Field]))
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
			} else {
				engine.log.Error("Error Get Element", sl.Err(err))
				return NO_TRANSITION
			}
		}
		engine.log.Debug("Got element", slog.String("elem", elem))
		engine.broadcast(common.Message{Type: common.ENGINE_ACTION, Ok: true, Body: map[string]any{
			"Action":       "GetElement",
			"Element":      elem,
			"LastElements": engine.Bag.LastElements(),
		}})
		return NO_TRANSITION
	}

}
