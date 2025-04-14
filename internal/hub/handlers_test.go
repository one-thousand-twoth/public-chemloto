package hub

// import (
// 	"bytes"
// 	"log/slog"
// 	"slices"
// 	"sync"
// 	"testing"

// 	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
// 	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
// 	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
// 	"github.com/gorilla/websocket"
// )

// var (
// 	buff       = &bytes.Buffer{}
// 	MockLogger = slog.New(slog.NewTextHandler(buff, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
// 		if a.Key == slog.TimeKey {
// 			return slog.Attr{}
// 		}
// 		return a
// 	}}))
// )

// func TestSubscribeToChannel(t *testing.T) {
// 	channels, hub, user := getDeps(t)
// 	t.Run("Subscribe User to channel", func(t *testing.T) {
// 		subscribeToChannel(SubscribeRequest{
// 			Type:   "HUB_SUBSCRIBE",
// 			Target: "room",
// 			Name:   "test_room",
// 		}, user.Name, hub,
// 		)
// 		cha, ok := channels.Get("test_room")
// 		if !ok {
// 			t.Fatal("Channel did not created")
// 		}
// 		if !slices.Contains(cha, user.conn) {
// 			t.Fatal("user not added to channel")
// 		}

// 	})

// }

// func TestSubscribeToRoom(t *testing.T) {
// 	channels, hub, user := getDeps(t)

// 	t.Run("SubscribeUser to Room", func(t *testing.T) {
// 		err := subscribeToRoom(hub, SubscribeRequest{
// 			Type:   "HUB_SUBSCRIBE",
// 			Target: "room",
// 			Name:   "test_room",
// 		},
// 			user.Name,
// 			MockLogger,
// 		)
// 		if err != nil {
// 			t.Errorf("Subsribe func return error: %s", err.Error())
// 		}
// 		// t.Log(channels.state)
// 		cha, ok := channels.Get("test_room")
// 		if !ok {
// 			t.Errorf("Channel did not created")
// 		}
// 		if !slices.Contains(cha, user.conn) {
// 			t.Errorf("user not added to channel")
// 		}
// 	})
// 	tests := []struct {
// 		name      string
// 		request   SubscribeRequest
// 		expectErr enerr.Kind
// 	}{
// 		{
// 			name: "SubscribeUser to Room (expected notExist error)",
// 			request: SubscribeRequest{
// 				Type:   "HUB_SUBSCRIBE",
// 				Target: "room",
// 				Name:   "test_room_failed",
// 			},
// 			expectErr: enerr.NotExist,
// 		},
// 		{
// 			name: "SubscribeUser (Target undefined)",
// 			request: SubscribeRequest{
// 				Type:   "HUB_SUBSCRIBE",
// 				Target: "",
// 				Name:   "test_room",
// 			},
// 			expectErr: enerr.Validation,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := subscribeToRoom(hub, tt.request, user.Name, MockLogger)
// 			if err == nil {
// 				t.Errorf("func expect error")
// 			}
// 			if !enerr.KindIs(tt.expectErr, err) {
// 				t.Errorf("expect %v error, got %v", tt.expectErr, err)
// 			}
// 		})
// 	}
// }

// func getDeps(t *testing.T) (*repository.ChannelsState, *Hub, *User) {
// 	channels := repository.NewChannelState()

// 	connectionsSt := &connectionsState{state: make(map[string]*SockConnection)}

// 	hub := &Hub{
// 		upgrader:      websocket.Upgrader{},
// 		log:           MockLogger,
// 		Rooms:         &roomsState{state: make(map[string]*Room)},
// 		Users:         &usersState{state: make(map[string]*User)},
// 		Connections:   connectionsSt,
// 		Channels:      channels,
// 		eventHandlers: make(map[string]HandlerFunc),
// 		eventChan:     make(chan internalEventWrap, 10),
// 	}
// 	err := hub.AddNewRoom(CreateRoomRequest{
// 		Name:        "test_room",
// 		MaxPlayers:  2,
// 		Elements:    map[string]int{},
// 		Time:        0,
// 		IsAuto:      false,
// 		IsAutoCheck: false,
// 	})
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	username := "Test User"
// 	user := &User{
// 		Name:     username,
// 		Apikey:   "apikey",
// 		Room:     "uuid",
// 		Role:     common.Player_Role,
// 		conn:     "conn_id",
// 		channels: []string{},
// 		mutex:    sync.Mutex{},
// 	}
// 	connectionsSt.state[user.conn] = &SockConnection{}
// 	connectionsSt.state[user.conn].MessageChan = make(chan common.Message, 20)

// 	hub.Users.Add(user)
// 	return channels, hub, user
// }
