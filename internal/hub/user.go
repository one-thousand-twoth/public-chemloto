package hub

import (
	"encoding/json"
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

type Role int

//go:generate stringer -type=Role -output=./stringer_types.go
const (
	NONE Role = iota
	Admin_Role
	Judge_Role
	Player_Role
)

type User struct {
	Name     string `json:"name"`
	Apikey   string `json:"token"`
	Room     string `json:"room"`
	Role     Role   `json:"role"`
	conn     string
	channels []string
	mutex    sync.Mutex
}

func (r *User) MarshalJSON() ([]byte, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	user := struct {
		Name string `json:"username"`
		Room string `json:"room"`
		Role string `json:"role"`
	}{r.Name, r.Room, r.Role.String()}
	return json.Marshal(user)
}

func NewUser(name string, apikey string, conn string, role Role, channels []string) *User {
	return &User{
		Name:     name,
		Apikey:   apikey,
		conn:     conn,
		channels: channels,
		Role:     role,
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
func (r *User) SetChannels(channels ...string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.channels = append(r.channels, channels...)
}

type usersState struct {
	// map key is username
	state map[string]*User
	mutex sync.RWMutex
}

func (rs *usersState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
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
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

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
