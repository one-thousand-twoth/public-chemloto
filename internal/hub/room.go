package hub

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	enmodels "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
)

type Engine interface {
	// Получить текущее состояние, например при перезагрузке страницы
	PreHook() map[string]any
	// Обработать событие
	Input(enmodels.Action)
	Start()
	AddPlayer(enmodels.Player) error
	RemovePlayer(name string) error
}

func (h *Hub) AddNewRoom(r *room) error {
	if err := h.Rooms.add(r); err != nil {
		return err
	}
	// Set Init function for all room channels
	h.Channels.SetChannelFunc(r.name, func(ch chan common.Message) {
		body := r.engine.PreHook()
		// h.log.Error("{Unsafe} sending engine state", "body", body)
		ch <- common.Message{
			Type: common.ENGINE_INFO,
			Ok:   true,
			Body: body,
		}
	})
	return nil
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
	// map key is a room name
	state map[string]*room
	mutex sync.RWMutex
}

func (rs *roomsState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
}

func (rs *roomsState) get(id string) (r *room, ok bool) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	r, ok = rs.state[id]
	return
}
func (rs *roomsState) add(room *room) error {
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
