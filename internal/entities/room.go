package entities

type Room struct {
	Name   string         `json:"name" validate:"required,min=1,safeinput"`
	Engine ExternalEngine `json:"engine"`
}

type ExternalEngine struct {
	Name       string
	Status     string
	Players    int
	MaxPlayers int
}
