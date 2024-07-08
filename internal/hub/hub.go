// hub отвечает за контроль над вебсокет клиентами
package hub

import (
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
)

type Hub struct {
	Rooms *roomsState
	Users *usersState

	// registerChan   chan Instruction
	// unregisterChan chan Instruction
	// broadcastChan  chan Message
	// castChan       chan CastMessage
	// log            *slog.Logger
	// storage        *sqlite.Storage
	// rooms          map[string]*Room
	// roomMutex      sync.RWMutex
	// freeUsers      map[string]*Client
	// usersMutex     sync.RWMutex
}

func NewHub(storage *sqlite.Storage, log *slog.Logger) *Hub {
	return &Hub{
		Rooms: &roomsState{state: make(map[string]*room)},
		Users: &usersState{state: make(map[string]*Client)},
	}
}
