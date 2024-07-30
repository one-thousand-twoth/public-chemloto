package hub

import (
	"encoding/json"
	"fmt"

	"github.com/anrew1002/Tournament-ChemLoto/internal/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type HandlerFunc func(*Hub, internalEventWrap)

func (f HandlerFunc) HandleEvent(h *Hub, msg internalEventWrap) {
	f(h, msg)
}

func (h *Hub) SetupHandlers() {
	h.UseHandler(HUB_SUBSCRIBE, Subscribe)
}

func (h *Hub) UseHandler(t MessageType, f HandlerFunc) {
	h.eventHandlers[t.String()] = f
}

// Subscribe обрабатывает логику добавления user`a в подписчики канала
func Subscribe(h *Hub, e internalEventWrap) {
	type dataT struct {
		Type   string
		Target string
		Name   string
	}
	op := "Subscribe handler"
	log := h.log.With("op", op)
	log.Debug("Start Handle Event", "usr", e.userId, "room", e.room, "data", e.msg)
	var data dataT
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		log.Error("failed to decode event body", sl.Err(err))
		return
	}
	if data.Target == "" {
		log.Error("empty field")
		return
	}
	if data.Name == "" {
		log.Error("empty field")
		return
	}

	usr, ok := h.Users.Get(e.userId)
	if !ok {
		log.Error("Error getting usr")
	}
	log.Debug("Context", "user", usr, "data", data)
	conn := usr.conn
	switch data.Target {
	case "room":
		oldRoom := usr.SetRoom(data.Name)
		h.Channels.Remove(oldRoom, conn)
		h.Channels.Add(data.Name, conn)
	case "channel":
		usr.SetChannels(data.Name)
		h.Channels.Add(data.Name, conn)
	default:
		log.Error(fmt.Sprintf("unknown target %s", data.Target))
		return
	}
	channelSubs, ok := h.Channels.Get(data.Name)
	if !ok {
		log.Error("Cant find channel")
		return
	}
	log.Debug("current status", "channelSubs", channelSubs, "user", conn)
	envelope, err := json.Marshal(data)
	if err != nil {
		log.Error("Marshaling data", "data", data)
	}
	h.SendMessageOverChannel(data.Name, models.Message{
		Type: websocket.TextMessage,
		Body: envelope,
	})

}
