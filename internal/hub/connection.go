package hub

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SockConnection struct {
	ID           string
	Conn         *websocket.Conn
	User         string
	CloseChannel chan struct{}

	mutex       sync.RWMutex
	MessageChan chan models.Message
}

func NewConnection(conn *websocket.Conn, user string) *SockConnection {
	return &SockConnection{
		ID:           uuid.New().String(),
		Conn:         conn,
		User:         user,
		CloseChannel: make(chan struct{}),
		MessageChan:  make(chan models.Message),
	}
}

func (r *SockConnection) MarshalJSON() ([]byte, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	SockConnection := struct {
		ID   string `json:"ID"`
		Name string `json:"name"`
		// Name string `json:"Name,omitempty"`
	}{r.ID, r.User}
	return json.Marshal(SockConnection)
}

// connectionsState - map структура с RWmutex
type connectionsState struct {
	state map[string]*SockConnection
	mutex sync.RWMutex
}

func (rs *connectionsState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
}

func (rs *connectionsState) Get(id string) *SockConnection {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.state[id]
}

func (rs *connectionsState) Add(sockConn *SockConnection) error {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	_, ok := rs.state[sockConn.ID]
	if ok {
		return errors.New("already exist SockConnection")
	} else {
		rs.state[sockConn.ID] = sockConn
	}

	return nil
}

func (rs *connectionsState) Remove(id string) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	delete(rs.state, id)
}
