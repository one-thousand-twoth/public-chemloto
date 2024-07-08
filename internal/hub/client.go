package hub

import (
	"errors"
	"sync"
)

type Client struct {
	Name   string
	apikey string
}

type usersState struct {
	state map[string]*Client
	mutex sync.RWMutex
}

func (rs *usersState) Get(id string) *Client {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.state[id]
}
func (rs *usersState) Add(name string, token string) error {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	_, ok := rs.state[name]
	if ok {
		return errors.New("already exist user")
	} else {
		rs.state[token] = &Client{Name: name, apikey: token}
	}

	return nil
}
