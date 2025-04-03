package hub

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
)

type internalEventWrap struct {
	userId         string
	msgType        common.MessageType
	MessageChannel chan common.Message
	msg            map[string]interface{}
}

func NewEventWrap(
	userID string,
	// user entities.User,
	// room string,
	// role common.Role,
	msg map[string]interface{},
	msgType common.MessageType,
	msgChan chan common.Message,
) internalEventWrap {
	return internalEventWrap{
		userId:         userID,
		MessageChannel: msgChan,
		msgType:        msgType,
		msg:            msg,
	}
}
