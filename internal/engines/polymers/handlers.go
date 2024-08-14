package polymers

import (
	"errors"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
)

type HandlerFunc func(models.Action) stateInt

// func (engine *PolymersEngine) SetupHandlers() {
// 	engine.UseHandler("GetElement", engine.GetElement())
// }
// func (engine *PolymersEngine) UseHandler(eventName string, handler HandlerFunc) {
// 	engine.handlers[eventName] = handler
// }

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
