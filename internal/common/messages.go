package common

import "encoding/json"

type Message struct {
	Type   MessageType
	Ok     bool
	Errors []error
	Body   map[string]any
}

type MessageType int

func (m MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

//go:generate stringer -type=MessageType
const UNDEFINED MessageType = -1
const (
	HUB_SUBSCRIBE MessageType = iota + 1
	ENGINE_ACTION
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
