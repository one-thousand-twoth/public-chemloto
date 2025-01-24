package hub

import (
	"bytes"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
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

func TestSubscribe(t *testing.T) {
	hub := NewHub(MockLogger, websocket.Upgrader{})

	username := "Test User"
	user := NewUser(username, "apikey", "uuid", common.Player_Role, []string{})
	hub.Users.Add(user)
	type args struct {
		h *Hub
		e internalEventWrap
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "with empty fields should error ",
			args: args{h: hub, e: internalEventWrap{}},
		},
		{
			name: "New subscription",
			args: args{h: hub, e: internalEventWrap{
				msgType: common.HUB_SUBSCRIBE,
				msg: map[string]interface{}{
					"Type":   "HUB_SUBSCRIBE",
					"Target": "room",
					"Name":   "test_room"},
				room:   "",
				userId: username,
			}},
		},
		{
			name: "Old subscription",
			args: args{h: hub, e: internalEventWrap{
				msgType: common.HUB_SUBSCRIBE,
				msg: map[string]interface{}{
					"Type":   "HUB_SUBSCRIBE",
					"Target": "channel",
					"Name":   "admin"},
				room:   "",
				userId: username,
			}},
		},
	}
	for k, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Subscribe(tt.args.h, tt.args.e)
			if k == 0 {
				if s := buff.String(); !strings.Contains(s, `level=ERROR msg="empty field"`) {
					t.Errorf("Unexpected log: %s", s)
				}

			}
			if k == 1 {
				if s := buff.String(); strings.Contains(s, `ERROR`) {
					t.Errorf("1_Unexpected log: %s", s)
				}
				l, ok := hub.Channels.Get("test_room")
				if !ok {
					t.Errorf("Channel test_room wasnt created")
				}
				if len(l) != 1 && l[0] != "uuid" {
					t.Errorf("Connection wasnt added to channel`s list")
				}
				if !reflect.DeepEqual(user.channels, []string{"test_room"}) {
					t.Logf("%+v", user.channels)
					t.Errorf("Channel wasnt added to user list")
				}
			}
			if k == 2 {
				if s := buff.String(); strings.Contains(s, `ERROR`) {
					t.Errorf("1_Unexpected log: %s", s)
				}
				l, ok := hub.Channels.Get("admin")
				if !ok {
					t.Errorf("Channel admin wasnt created")
				}
				if len(l) != 1 && l[0] != "uuid" {
					t.Errorf("Connection wasnt added to channel`s list")
				}
				if !reflect.DeepEqual(user.channels, []string{"test_room", "admin"}) {
					t.Logf("%+v", user.channels)
					t.Errorf("Channel wasnt added to user list")
				}
			}
			buff.Reset()
		})
	}
}
func TestSubscribe2(t *testing.T) {
	channels := &channelsState{
		state: make(map[string]mapset.Set[string]),
		initFunctions: map[string]func(chan common.Message){
			"default": func(c chan common.Message) {},
		},
	}
	hub := &Hub{
		upgrader:      websocket.Upgrader{},
		log:           MockLogger,
		Rooms:         &roomsState{state: make(map[string]*Room)},
		Users:         &usersState{state: make(map[string]*User)},
		Connections:   &connectionsState{state: make(map[string]*SockConnection)},
		Channels:      channels,
		eventHandlers: make(map[string]HandlerFunc),
		eventChan:     make(chan internalEventWrap, 10),
	}

	username := "Test User"
	user := &User{
		Name:     username,
		Apikey:   "apikey",
		Room:     "uuid",
		Role:     common.Player_Role,
		conn:     "",
		channels: []string{},
		mutex:    sync.Mutex{},
	}

	hub.Users.Add(user)
	t.Run("Subscribe User to channel", func(t *testing.T) {
		event := internalEventWrap{
			msgType: common.HUB_SUBSCRIBE,
			msg: map[string]interface{}{
				"Type":   "HUB_SUBSCRIBE",
				"Target": "channel",
				"Name":   "test_channel"},
			room:   "",
			userId: username,
		}
		Subscribe(hub, event)

		if !channels.state["test_channel"].Contains(user.conn) {
			t.Fatal("user not added to channel")
		}

	})
}
