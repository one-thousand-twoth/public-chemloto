package engines

import (
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
)

type EngineBuilder struct {
}

func getEngine(engine string, config any) (any, error) {
	switch engine {
	case "polymers":
		return createPolymerEngine(config), nil
	}

	return nil, enerr.E("No such engine declared")
}

func createPolymerEngine(config any) any {
	return polymers.New(
		h.log.With(slog.String("room", r.Name)),
		polymers.PolymersEngineConfig{
			Elements:   r.Elements,
			Checks:     checks,
			TimerInt:   r.Time,
			MaxPlayers: r.MaxPlayers,
			Unicast: func(userID string, msg common.Message) {
				go func() {
					h.log.Debug("Unicast message")
					usr, ok := h.Users.Get(userID)
					if !ok {
						h.log.Error("failed to get user while Unicast message from engine")
						return
					}
					connID := usr.GetConnection()
					conn, ok := h.Connections.Get(connID)
					if !ok {
						h.log.Error("failed to get user connection while Unicast message from engine")
						return
					}
					conn.MessageChan <- msg
				}()
			},
			Broadcast: func(msg common.Message) {
				h.log.Debug("Broadcast message")
				go h.SendMessageOverChannel(r.Name, msg)
			},
		},
	)
}
