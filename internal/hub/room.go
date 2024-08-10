package hub

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/engine"
)

type Engine interface {
	// Получить текущее состояние, например при перезагрузке страницы
	PreHook()
	// Обработать событие
	Input(engine.Action)
}

type room struct {
	name       string
	maxPlayers int
	elements   map[string]int
	time       int
	isAuto     bool
	engine     Engine
}

func NewRoom(
	name string,
	maxPlayers int,
	elements map[string]int,
	time int,
	isAuto bool,
	engine Engine,
) *room {
	return &room{
		name:       name,
		maxPlayers: maxPlayers,
		elements:   elements,
		time:       time,
		isAuto:     isAuto,
		engine:     engine,
	}
}

func (r *room) MarshalJSON() ([]byte, error) {
	room := struct {
		Name   string         `json:"name"`
		Max    int            `json:"maxPlayers" `
		Elems  map[string]int `json:"elementCounts"`
		Time   int            `json:"time" `
		IsAuto bool           `json:"isAuto"`
	}{r.name, r.maxPlayers, r.elements, r.time, r.isAuto}
	return json.Marshal(room)
}

type roomsState struct {
	state map[string]*room
	mutex sync.RWMutex
}

func (rs *roomsState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
}

func (rs *roomsState) Get(id string) *room {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.state[id]
}
func (rs *roomsState) Add(room *room) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	_, ok := rs.state[room.name]
	if ok {
		return errors.New("already exist room")
	} else {
		rs.state[room.name] = room
	}

	return nil
}
