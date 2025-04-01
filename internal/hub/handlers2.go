package hub

import (
	"fmt"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/usecase"
	"github.com/mitchellh/mapstructure"
)

// Subscribe обрабатывает логику добавления user`a в подписчики канала
func SubscribeHandler2(h *Hub, e internalEventWrap) error {

	const op enerr.Op = "Subscribe handler"
	log := h.log.With("op", op)

	log.Debug("Start Handle Event", "usr", e.userId, "room", e.room, "data", fmt.Sprintf("%v", e.msg))

	var data SubscribeRequest
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		return enerr.E(op, fmt.Sprintf("failed to decode event body: %s", err.Error()))
	}
	// conn, ok := h.Connections.Get(e.connId)
	// if !ok {
	// 	return enerr.E(op, "Failed getting connection of user")
	// }

	switch data.Target {
	case "room":
		err := usecase.SubscribeToRoom(h.Rooms2, data.Name, &e.user)
		if err != nil {
			return err
		}
	case "channel":
		// err := usecase.SubscribeToChannel(h.Channels, data.Name, e.user)
		panic("TODO")
	default:
		return enerr.E(op, fmt.Sprintf("unknown target %s", data.Target))
	}

	e.MessageChannel <- common.Message{
		Type: common.HUB_SUBSCRIBE,
		Ok:   true,
		Body: map[string]any{
			"Target": "room",
			"Name":   data.Name,
		},
	}
	fun, ok := h.Channels.GetChannelFunc(data.Name)
	if ok {
		fun(e.MessageChannel)
	}
	return nil
}
