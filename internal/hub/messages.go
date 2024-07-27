package hub

type MessageType int

//go:generate stringer -type=MessageType
const UNDEFINED MessageType = -1
const (
	HUB_SUBSCRIBE MessageType = iota + 1
	HUB_NEW_ROOM
)

var MapEnumStringToMessageType = func() map[string]MessageType {
	m := make(map[string]MessageType)
	for i := HUB_SUBSCRIBE; i <= HUB_NEW_ROOM; i++ {
		m[i.String()] = i
	}
	return m
}()

type internalEventWrap struct {
	userId  string
	room    string
	msgType MessageType
	msg     map[string]interface{}
}

func NewEventWrap(userID string, room string, msg map[string]interface{}, msgType MessageType) internalEventWrap {
	return internalEventWrap{
		userId:  userID,
		room:    room,
		msgType: msgType,
		msg:     msg,
	}
}
