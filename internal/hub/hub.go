// hub отвечает за контроль над вебсокет клиентами
package hub

import (
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
	mapset "github.com/deckarep/golang-set/v2"
)

type Hub struct {
	Rooms       *roomsState
	Users       *usersState
	Connections *connectionsState
	Channels    *channelsState

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
		Rooms:       &roomsState{state: make(map[string]*room)},
		Users:       &usersState{state: make(map[string]*User)},
		Connections: &connectionsState{state: make(map[string]*SockConnection)},
		Channels:    &channelsState{state: make(map[string]mapset.Set[string])},
	}
}

func RemoveComparable[T comparable, V comparable](m map[T]V, element V) map[T]V {
	// var result []T
	for k, item := range m {
		if item == element {
			delete(m, k)
		}
	}
	return m
}
func (h *Hub) ListenClient(connection *SockConnection) {

}

func (h *Hub) Input(input string) {

}
