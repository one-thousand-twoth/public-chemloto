package hub

import (
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	enmodels "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
)

type HandlerFunc func(*Hub, internalEventWrap)

func (f HandlerFunc) HandleEvent(h *Hub, msg internalEventWrap) {
	f(h, msg)
}

func (h *Hub) SetupHandlers() {
	h.UseHandler(common.HUB_SUBSCRIBE, Subscribe)
	h.UseHandler(common.ENGINE_ACTION, EngineAction)
	h.UseHandler(common.HUB_STARTGAME, StartGame)
}

func (h *Hub) UseHandler(t common.MessageType, f HandlerFunc) {
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
	log.Debug("Start Handle Event", "usr", e.userId, "room", e.room, "data", fmt.Sprintf("%v", e.msg))
	var data dataT
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		log.Error("failed to decode event body", sl.Err(err))
		return
	}
	if data.Target == "" || data.Name == "" {
		log.Error("empty field")
		return
	}

	usr, ok := h.Users.Get(e.userId)
	if !ok {
		log.Error("Error getting usr")
		return
	}
	usr.mutex.Lock()
	defer usr.mutex.Unlock()

	log.Debug("Context for Subscribe Event", "user", usr.Name, "data", data)
	connID := usr.conn
	conn, ok := h.Connections.Get(connID)
	if !ok {
		log.Error("Failed getting connection of user")
		return
	}
	switch data.Target {
	case "room":
		room, ok := h.Rooms.get(data.Name)
		if !ok {
			conn.MessageChan <- common.Message{
				Type:   common.HUB_SUBSCRIBE,
				Ok:     false,
				Errors: []string{"Не существующая комната"},
			}
			return
		}
		if err := room.engine.AddPlayer(e.userId); err != nil {
			conn.MessageChan <- common.Message{
				Type:   common.HUB_SUBSCRIBE,
				Ok:     false,
				Errors: []string{"Комната уже запущена, неё нельзя войти"},
			}
			return
		}

		oldRoomID := usr.SetRoom(data.Name)
		oldRoom, ok := h.Rooms.get(oldRoomID)
		if ok {
			if err := oldRoom.engine.RemovePlayer(e.userId); err != nil {
				log.Error("Failed to delete player", sl.Err(err))
				return
			}
		} else {
			log.Error("Failed to get room")
		}
		h.Channels.Remove(oldRoomID, connID)
		h.Channels.Add(data.Name, connID)
		// h.Channels.
	case "channel":
		usr.channels = append(usr.channels, data.Name)
		h.Channels.Add(data.Name, connID)
	default:
		log.Error(fmt.Sprintf("unknown target %s", data.Target))
		return
	}
	channelSubs, ok := h.Channels.Get(data.Name)
	if !ok {
		log.Error("Cant find channel")
		return
	}
	log.Debug("current status", "channelSubs", channelSubs, slog.Any("user", usr))

	conn.MessageChan <- common.Message{
		Type: common.HUB_SUBSCRIBE,
		Ok:   true,
		Body: map[string]any{
			"Target": "room",
			"Name":   data.Name,
		},
	}
	fun, ok := h.Channels.GetChannelFunc(data.Name)
	if !ok {
		return
	}
	fun(conn.MessageChan)
}
func EngineAction(h *Hub, e internalEventWrap) {
	room, ok := h.Rooms.get(e.room)
	if !ok {
		h.log.Error("Cannot find room for EngineEvent", "room", e.room)
		return
	}
	room.engine.Input(enmodels.Action{
		Player:   e.userId,
		Envelope: e.msg,
	})
}

func StartGame(h *Hub, e internalEventWrap) {
	type dataT struct {
		Type string
		Name string
	}
	op := "StartGame handler"
	log := h.log.With("op", op)
	log.Debug("Start Handle Event", "usr", e.userId, "room", e.room, "data", fmt.Sprintf("%v", e.msg))

	room, ok := h.Rooms.get(e.room)
	if !ok {
		log.Error("Failed getting room")
		return
	}
	room.engine.Start()
}
