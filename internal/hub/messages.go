package hub

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
)

type internalEventWrap struct {
	userId string
	// connId  string
	user           entities.User
	room           string
	role           common.Role
	msgType        common.MessageType
	MessageChannel chan common.Message
	msg            map[string]interface{}
}

func NewEventWrap(
	userID string,
	user entities.User,
	room string,
	role common.Role,
	msg map[string]interface{},
	msgType common.MessageType,
	msgChan chan common.Message,
) internalEventWrap {
	return internalEventWrap{
		userId:         userID,
		user:           user,
		MessageChannel: msgChan,
		room:           room,
		role:           role,
		msgType:        msgType,
		msg:            msg,
	}
}
