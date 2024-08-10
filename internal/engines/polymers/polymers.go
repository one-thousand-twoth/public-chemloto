package polymers

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/engine"
)

func New(log *slog.Logger) *PolymersEngine {
	// var src cryptorand.CryptoSource
	src := rand.NewSource(time.Now().UnixNano())
	return &PolymersEngine{
		log:              log.With(slog.String("source", "PulymersEngine")),
		lastElementsKeys: []string{},
		Elements:         map[string]int{},
		pushedElements:   []string{},
		ActionChan:       make(chan engine.Action),
		handlers:         map[string]HandlerFunc{},
		rnd:              rand.New(src),
	}
}

type Event struct {
	Action string
}

type PolymersEngine struct {
	log *slog.Logger
	// Названия элементов в игре
	lastElementsKeys []string
	// Общий мешок
	Elements map[string]int
	// Элементы которые достали из мешка
	pushedElements []string
	// Обработка действий игроков
	ActionChan chan engine.Action
	handlers   map[string]HandlerFunc

	rnd *rand.Rand
}

func (engine *PolymersEngine) Start() {
	go func() {
		for e := range engine.ActionChan {
			action, ok := e.Envelope["Action"].(string)
			if !ok {
				engine.log.Error("failed to extract Action field")
			}
			handle, ok := engine.handlers[action]
			if !ok {
				engine.log.Error(fmt.Sprintf("No handler for event %s", action))
				continue
			}
			engine.log.Debug("Start Handling Engine Event")
			handle(e)
		}
	}()
}

func (engine *PolymersEngine) Input(e engine.Action) {
	engine.ActionChan <- e
}
func (engine *PolymersEngine) PreHook() {

}
