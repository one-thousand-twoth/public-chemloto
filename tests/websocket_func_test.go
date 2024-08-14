package tests

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/gavv/httpexpect/v2"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebsocketSubscribe(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	username := getRandomUsername()
	resp := Createuser(e, username, "")
	uWS := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   baseUrl,
	}
	conn, _, err := websocket.DefaultDialer.Dial(uWS.String()+"/ws"+"?token="+resp.Value("token").String().Raw(), nil)
	assert.Nil(t, err, "Websocket connection cannot be establish")
	// watcher, _, err := websocket.DefaultDialer.Dial(uWS.String()+"/ws"+"?token="+resp.Value("token").String().Raw(), nil)
	type dataT struct {
		Type   string
		Target string
		Name   string
	}
	data := dataT{Type: common.HUB_SUBSCRIBE.String(), Target: "room", Name: "roomname"}
	conn.WriteJSON(data)

	messageType, msg, err := conn.ReadMessage()
	assert.Nil(t, err, "Ошибка при чтении сообщения")
	assert.Equal(t, messageType, websocket.TextMessage)
	t.Log(string(msg))
	var output dataT
	json.Unmarshal(msg, &output)
	if !reflect.DeepEqual(data, output) {
		t.Error("Не пришло подтверждение  о подписке")
	}
}
