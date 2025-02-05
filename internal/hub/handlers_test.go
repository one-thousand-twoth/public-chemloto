package hub

import (
	"bytes"
	"log/slog"
	"sync"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
)

var (
	buff       = &bytes.Buffer{}
	MockLogger = slog.New(slog.NewTextHandler(buff, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	}}))
)

type StubUserStore struct {
}

func (ps StubUserStore) Get(name string) (*User, bool) {
	return nil, false
}
func (ps StubUserStore) GetByToken(token string) (*User, bool) {
	return nil, false
}
func (ps StubUserStore) Add(user *User) error {
	return nil
}
func (ps StubUserStore) Remove(username string) {

}

type StubChannelsStore struct {
	// map key is a channel name
	state map[string]mapset.Set[string]
	// Функции каналов срабатывающие при первом подключении.
	// key: channel name
	initFunctions map[string]func(chan common.Message)
}

func (cs StubChannelsStore) Get(channel string) ([]string, bool) {

	return cs.state[channel].ToSlice(), true
}
func (cs StubChannelsStore) Add(channel string, connection string) {

}
func (cs StubChannelsStore) Remove(channel string, connection string) {

}
func (cs StubChannelsStore) GetChannelFunc(channel string) (func(chan common.Message), bool) {
	return nil, false
}
func (cs StubChannelsStore) SetChannelFunc(channel string, fun func(chan common.Message)) {

}

func TestSubscribeInternal(t *testing.T) {
	channels := &channelsState{
		state: make(map[string]mapset.Set[string]),
		initFunctions: map[string]func(chan common.Message){
			"default": func(c chan common.Message) {},
		},
	}

	connectionsSt := &connectionsState{state: make(map[string]*SockConnection)}

	hub := &Hub{
		upgrader:      websocket.Upgrader{},
		log:           MockLogger,
		Rooms:         &roomsState{state: make(map[string]*Room)},
		Users:         &usersState{state: make(map[string]*User)},
		Connections:   connectionsSt,
		Channels:      channels,
		eventHandlers: make(map[string]HandlerFunc),
		eventChan:     make(chan internalEventWrap, 10),
	}
	err := hub.AddNewRoom(CreateRoomRequest{
		Name:        "test_room",
		MaxPlayers:  2,
		Elements:    map[string]int{},
		Time:        0,
		IsAuto:      false,
		IsAutoCheck: false,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	username := "Test User"
	user := &User{
		Name:     username,
		Apikey:   "apikey",
		Room:     "uuid",
		Role:     common.Player_Role,
		conn:     "conn_id",
		channels: []string{},
		mutex:    sync.Mutex{},
	}
	connectionsSt.state[user.conn] = &SockConnection{}
	connectionsSt.state[user.conn].MessageChan = make(chan common.Message, 20)

	hub.Users.Add(user)
	t.Run("Subscribe User to channel", func(t *testing.T) {
		event := internalEventWrap{
			msgType: common.HUB_SUBSCRIBE,
			msg: map[string]interface{}{
				"Type":   "HUB_SUBSCRIBE",
				"Target": "channel",
				"Name":   "default"},
			room:   "",
			userId: username,
		}
		err := SubscribeHandler(hub, event)
		if err != nil {
			t.Fatalf("Subsribe func return error: %s", err.Error())
		}
		t.Log(channels.state)
		cha, ok := channels.state["default"]
		if !ok {
			t.Fatal("Channel did not created")
		}
		if !cha.Contains(user.conn) {
			t.Fatal("user not added to channel")
		}

	})

	t.Run("SubscribeUser to Room", func(t *testing.T) {
		event := internalEventWrap{
			msgType: common.HUB_SUBSCRIBE,
			msg: map[string]interface{}{
				"Type":   "HUB_SUBSCRIBE",
				"Target": "room",
				"Name":   "test_room"},
			room:   "",
			userId: username,
		}
		err := SubscribeHandler(hub, event)
		if err != nil {
			t.Fatalf("Subsribe func return error: %s", err.Error())
		}
		t.Log(channels.state)
		cha, ok := channels.state["test_room"]
		if !ok {
			t.Fatal("Channel did not created")
		}
		if !cha.Contains(user.conn) {
			t.Fatal("user not added to channel")
		}
	})
}

func Test_subscribeToRoom(t *testing.T) {
	channels := &channelsState{
		state: make(map[string]mapset.Set[string]),
		initFunctions: map[string]func(chan common.Message){
			"default": func(c chan common.Message) {},
		},
	}

	connectionsSt := &connectionsState{state: make(map[string]*SockConnection)}

	hub := &Hub{
		upgrader:      websocket.Upgrader{},
		log:           MockLogger,
		Rooms:         &roomsState{state: make(map[string]*Room)},
		Users:         &usersState{state: make(map[string]*User)},
		Connections:   connectionsSt,
		Channels:      channels,
		eventHandlers: make(map[string]HandlerFunc),
		eventChan:     make(chan internalEventWrap, 10),
	}
	err := hub.AddNewRoom(CreateRoomRequest{
		Name:        "test_room",
		MaxPlayers:  2,
		Elements:    map[string]int{},
		Time:        0,
		IsAuto:      false,
		IsAutoCheck: false,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	username := "Test User"
	user := &User{
		Name:     username,
		Apikey:   "apikey",
		Room:     "uuid",
		Role:     common.Player_Role,
		conn:     "conn_id",
		channels: []string{},
		mutex:    sync.Mutex{},
	}
	connectionsSt.state[user.conn] = &SockConnection{}
	connectionsSt.state[user.conn].MessageChan = make(chan common.Message, 20)

	hub.Users.Add(user)

	t.Run("SubscribeUser to Room", func(t *testing.T) {
		err := subscribeToRoom(hub, SubscribeRequest{
			Type:   "HUB_SUBSCRIBE",
			Target: "room",
			Name:   "test_room",
		},
			user,
			MockLogger,
			user.conn,
		)
		if err != nil {
			t.Fatalf("Subsribe func return error: %s", err.Error())
		}
		t.Log(channels.state)
		cha, ok := channels.state["test_room"]
		if !ok {
			t.Fatal("Channel did not created")
		}
		if !cha.Contains(user.conn) {
			t.Fatal("user not added to channel")
		}
	})
	t.Run("SubscribeUser to Room (expected notExist errror)", func(t *testing.T) {
		err := subscribeToRoom(hub, SubscribeRequest{
			Type:   "HUB_SUBSCRIBE",
			Target: "room",
			Name:   "test_room_failed",
		},
			user,
			MockLogger,
			user.conn,
		)
		t.Log(hub.Rooms.state)
		if err == nil {
			t.Fatalf("func expect error ")
		}
		if !enerr.KindIs(enerr.NotExist, err) {
			t.Fatal("expect empty field error")
		}
	})

	t.Run("SubscribeUser (Target undefined)", func(t *testing.T) {
		err := subscribeToRoom(hub, SubscribeRequest{
			Type:   "HUB_SUBSCRIBE",
			Target: "",
			Name:   "test_room",
		},
			user,
			MockLogger,
			user.conn,
		)
		if err == nil {
			t.Fatalf("Subsribe func expect error: ")
		}
		if !enerr.KindIs(enerr.Validation, err) {
			t.Fatal("expect empty field error")
		}
	})
}
