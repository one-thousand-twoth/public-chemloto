package entities

type Room struct {
	Name   string         `json:"name"`
	Engine ExternalEngine `json:"engine"`
}

type ExternalEngine struct {
	Name       string
	Status     string
	Players    int
	MaxPlayers int
}
