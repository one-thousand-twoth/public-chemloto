package engines

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	enmodels "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
	"github.com/anrew1002/Tournament-ChemLoto/web"
	"github.com/mitchellh/mapstructure"
)

type EngineBuilder struct {
}

func NewEngine(
	engineType string,
	name string,
	log *slog.Logger,
	config map[string]interface{},
) (enmodels.Engine, error) {
	const op enerr.Op = "engines.factory/GetEngine"
	log = log.With(slog.String("room", name))
	switch engineType {
	case "polymers":
		var data PolymersConfig
		if err := mapstructure.Decode(config, &data); err != nil {
			return nil, enerr.E(op, err, enerr.InvalidRequest)
		}
		return createPolymerEngine(log, data), nil
	}

	return nil, enerr.E("No such engine declared")
}

type PolymersConfig struct {
	Name        string         `json:"name" validate:"required,min=1,safeinput"`
	MaxPlayers  int            `json:"maxPlayers" validate:"required,gt=1,lt=100"`
	Elements    map[string]int `json:"elementCounts" validate:"required"`
	Time        int            `validate:"excluded_if=isAuto false,gte=0"`
	IsAuto      bool           `json:"isAuto"`
	IsAutoCheck bool           `json:isAutoCheck`
}

// Returns the structure needed to check the collected chem-structures of the player
//
// NOTE: нужно будет сделать чтобы парсился файл один раз при запуске проекта,
// и падал с ошибкой тоже только при запуске, а не в процессе.
func parseEngineJson() polymers.Checks {
	jsonFile, err := web.Polymers.Open("polymers.json")
	if err != nil {
		panic(fmt.Sprintf("while parsing polymers json: %v", err))
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Sprintf("while parsing polymers json: %v", err))
	}
	var checks polymers.Checks
	err = json.Unmarshal(byteValue, &checks.Fields)
	if err != nil {
		panic(fmt.Sprintf("while parsing polymers json: %v", err))
	}
	return checks
}

func createPolymerEngine(log *slog.Logger, config PolymersConfig) *polymers.PolymersEngine {
	checks := parseEngineJson()
	return polymers.New(
		log,
		polymers.PolymersEngineConfig{
			Elements:   config.Elements,
			Checks:     checks,
			TimerInt:   config.Time,
			MaxPlayers: config.MaxPlayers,
			Unicast: func(userID string, msg common.Message) {
				go func() {
					log.Debug("Unicast message")
					log.Error("Unimplemented")
					// usr, ok := h.Users.Get(userID)
					// if !ok {
					// 	log.Error("failed to get user while Unicast message from engine")
					// 	return
					// }
					// connID := usr.GetConnection()
					// conn, ok := h.Connections.Get(connID)
					// if !ok {
					// 	h.log.Error("failed to get user connection while Unicast message from engine")
					// 	return
					// }
					// conn.MessageChan <- msg
				}()
			},
			Broadcast: func(msg common.Message) {
				log.Debug("Broadcast message")
				log.Error("Unimplemented")
				// go h.SendMessageOverChannel(config.Name, msg)
			},
		},
	)
}
