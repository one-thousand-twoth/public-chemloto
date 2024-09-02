package common

import (
	"encoding/json"
	"errors"
)

type Message struct {
	Type   MessageType
	Ok     bool
	Errors []string
	Body   map[string]any
}

type MessageType int

func (m MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
func (m *MessageType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	inType, ok := MapEnumStringToMessageType[str]
	if !ok {
		return errors.New("unknown messageType")
	}
	// fmt.Println(inType)
	*m = inType
	return nil
}

//go:generate stringer -type=MessageType
const UNDEFINED MessageType = -1
const (
	HUB_SUBSCRIBE MessageType = iota + 1
	HUB_UNSUBSCRIBE
	ENGINE_ACTION
	ENGINE_INFO
	HUB_STARTGAME
	// HUB_NEW_ROOM указывается последним чтобы работала функция [MapEnumStringToMessageType]
	HUB_NEW_ROOM
)

var MapEnumStringToMessageType = func() map[string]MessageType {
	m := make(map[string]MessageType)
	for i := HUB_SUBSCRIBE; i <= HUB_NEW_ROOM; i++ {
		m[i.String()] = i
	}
	return m
}()
