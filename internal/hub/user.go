package hub

import (
	"errors"
	"sync"
)

type User struct {
	Name     string
	apikey   string
	conn     *SockConnection
	channels []string
	mutex    sync.Mutex
}

func NewUser(name string, apikey string, conn *SockConnection, channels []string) *User {
	return &User{
		Name:     name,
		apikey:   apikey,
		conn:     conn,
		channels: channels,
	}
}

func (r *User) GetChannels() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.channels
}

func (r *User) SetConnection(conn *SockConnection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.conn = conn
}

type usersState struct {
	state map[string]*User
	mutex sync.RWMutex
}

func (rs *usersState) Get(id string) *User {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.state[id]
}

// TODO: Не работает для обнаружения конфликтов
func (rs *usersState) Add(user *User) error {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	for _, usr := range rs.state {
		if usr.Name == user.Name {
			return errors.New("already exists")
		}
	}
	rs.state[user.apikey] = user
	return nil
}

// func (users *usersState) Remove(user string, connection string) {
// 	users.mutex.Lock()
// 	defer users.mutex.Unlock()

// 	use
// 	// if usersExists {
// 	// 	users.state[user] = common.RemoveString(usersEntry, connection)
// 	// }
// }
