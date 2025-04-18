package tests

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/gavv/httpexpect/v2"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebsocketSubscribe(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	username := getRandomUsername()
	resp := Createuser(e, username, "").JSON().Object()
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
	var output common.Message
	err = json.Unmarshal(msg, &output)
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
	assert.Equal(t, common.HUB_SUBSCRIBE, output.Type)
	assert.Equal(t, output.Ok, true)

	val := httpexpect.NewValue(t, output.Body)
	val.Object().Value("Name").IsEqual("roomname")
	val.Object().Value("Target").IsEqual("room")
}
