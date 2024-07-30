package hub

import (
	"bytes"
	"log/slog"
	"reflect"
	"strings"
	"testing"
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

func TestSubscribe(t *testing.T) {
	hub := NewHub(nil, MockLogger)
	username := "Test User"
	user := NewUser(username, "apikey", "uuid", Player_Role, []string{})
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
				msgType: HUB_SUBSCRIBE,
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
				msgType: HUB_SUBSCRIBE,
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
