package hub

import "github.com/anrew1002/Tournament-ChemLoto/internal/common"

type internalEventWrap struct {
	userId  string
	room    string
	msgType common.MessageType
	msg     map[string]interface{}
}

func NewEventWrap(userID string, room string, msg map[string]interface{}, msgType common.MessageType) internalEventWrap {
	return internalEventWrap{
		userId:  userID,
		room:    room,
		msgType: msgType,
		msg:     msg,
	}
}
