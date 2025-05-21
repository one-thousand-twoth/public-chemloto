package engines

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

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
	config map[string]any,
	unicast enmodels.UnicastFunction,
	broadcast enmodels.BroadcastFunction,
) (enmodels.Engine, error) {
	const op enerr.Op = "engines.factory/GetEngine"
	switch engineType {
	case "polymers":
		var data PolymersConfig
		if err := mapstructure.Decode(config, &data); err != nil {
			return nil, enerr.E(op, err, enerr.InvalidRequest)
		}
		return createPolymerEngine(log, data, unicast, broadcast), nil
	}

	return nil, enerr.E("No such engine declared")
}

type PolymersConfig struct {
	Name        string         `mapstructure:"name"`
	MaxPlayers  int            `mapstructure:"maxPlayers"`
	Elements    map[string]int `mapstructure:"elementCounts"`
	Time        int            `mapstructure:"time"`
	IsAuto      bool           `mapstructure:"isAuto"`
	IsAutoCheck bool           `mapstructure:isAutoCheck`
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

func createPolymerEngine(
	log *slog.Logger,
	config PolymersConfig,
	unicast enmodels.UnicastFunction,
	broadcast enmodels.BroadcastFunction,
) *polymers.PolymersEngine {
	checks := parseEngineJson()
	return polymers.New(
		log,
		polymers.PolymersEngineConfig{
			Elements:   config.Elements,
			Checks:     checks,
			TimerInt:   config.Time,
			MaxPlayers: config.MaxPlayers,
			Unicast:    unicast,
			Broadcast:  broadcast,
		},
	)
}
