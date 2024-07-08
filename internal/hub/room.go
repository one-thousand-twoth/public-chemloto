package hub

import (
	"encoding/json"
	"errors"
	"sync"
)

type room struct {
	// id   string
	name string
}

func (r *room) MarshalJSON() ([]byte, error) {
	room := struct {
		Name string `json:"name"`
		// Name string `json:"Name,omitempty"`
	}{r.name}
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
func (rs *roomsState) Add(name string) error {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	_, ok := rs.state[name]
	if ok {
		return errors.New("already exist room")
	} else {
		rs.state[name] = &room{name: name}
	}

	return nil
}
