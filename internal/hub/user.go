package hub

import (
	"errors"
	"sync"
)

func remove[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

type User struct {
	Name     string `json:"name"`
	Apikey   string `json:"token"`
	Room     string `json:room`
	conn     string
	channels []string
	mutex    sync.Mutex
}

func NewUser(name string, apikey string, conn string, channels []string) *User {
	return &User{
		Name:     name,
		Apikey:   apikey,
		conn:     conn,
		channels: channels,
	}
}

func (r *User) GetChannels() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.channels
}

func (r *User) SetConnection(conn string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.conn = conn
}
func (r *User) SetRoom(room string) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	oldRoom := r.Room
	if oldRoom == room {
		return ""
	}
	remove(r.channels, oldRoom)
	r.channels = append(r.channels, room)
	r.Room = room
	return oldRoom
}
func (r *User) SetChannels(channel string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.channels = append(r.channels, channel)
}

type usersState struct {
	// map key is username
	state map[string]*User
	mutex sync.RWMutex
}

func (rs *usersState) Get(name string) (*User, bool) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	usr, ok := rs.state[name]
	return usr, ok
}

func (rs *usersState) GetByToken(token string) (*User, bool) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	for _, usr := range rs.state {
		if usr.Apikey == token {
			return usr, true
		}
	}
	return nil, false
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
	rs.state[user.Name] = user
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
