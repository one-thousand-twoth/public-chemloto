package entities

import "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"

type Room struct {
	Name   string        `json:"name"`
	Engine models.Engine `json:"engine"`
}

// type ExternalEngine struct {
// 	Name       string
// 	Status     string
// 	Players    int
// 	MaxPlayers int
// }
