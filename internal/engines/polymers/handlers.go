package polymers

import (
	"errors"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/engine"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
)

type HandlerFunc func(engine.Action)

func (eng *PolymersEngine) SetupHandlers() {
	eng.UseHandler("GetElement", eng.GetElement())
}
func (eng *PolymersEngine) UseHandler(eventName string, handler HandlerFunc) {
	eng.handlers[eventName] = handler
}

func (eng *PolymersEngine) GetElement() HandlerFunc {
	return func(e engine.Action) {
		elem, err := eng.getRandomElement()
		if err != nil {
			if errors.Is(err, ErrEmptyBag) {
				eng.log.Info("Empty bag!")
				return
			}
			eng.log.Error("Error Get Element", sl.Err(err))
			return
		}
		eng.log.Debug("Got element", slog.String("elem", elem))
	}

}
