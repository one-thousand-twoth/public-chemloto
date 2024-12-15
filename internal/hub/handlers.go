package hub

import (
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	enmodels "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/mitchellh/mapstructure"
)

type HandlerFunc func(*Hub, internalEventWrap)

func (f HandlerFunc) HandleEvent(h *Hub, msg internalEventWrap) {
	f(h, msg)
}

func (h *Hub) SetupHandlers() {
	h.UseHandler(common.HUB_SUBSCRIBE, Subscribe)
	h.UseHandler(common.HUB_UNSUBSCRIBE, UnSubscribe)
	h.UseHandler(common.ENGINE_ACTION, EngineAction)
	h.UseHandler(common.HUB_STARTGAME, StartGame)
	h.UseHandler(common.HUB_EXITGAME, StopGame)

}

// use handler f for messages t
func (h *Hub) UseHandler(t common.MessageType, f HandlerFunc) {
	h.eventHandlers[t.String()] = f
}

func StopGame(h *Hub, e internalEventWrap) {
	type dataT struct {
		Type string
		Name string
	}
	op := "Hub/StopGame"
	log := h.log.With("op", op)
	log.Debug("Start Handle Event StopGame", "usr", e.userId, "room", e.room, "data", fmt.Sprintf("%v", e.msg))
	var data dataT
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		log.Error("failed to decode event body", sl.Err(err))
		return
	}

	usr, ok := h.Users.Get(e.userId)
	if !ok {
		log.Error("Error getting usr")
		return
	}
	if usr.Role == common.Player_Role {
		log.Error("No Access")
		return
	}

	log.Debug("Context for Subscribe Event", "user", usr.Name, "data", data)
	connID := usr.conn
	conn, ok := h.Connections.Get(connID)
	if !ok {
		log.Error("Failed getting connection of user")
		return
	}
	room, ok := h.Rooms.get(data.Name)
	if !ok {
		conn.MessageChan <- common.Message{
			Type:   common.HUB_SUBSCRIBE,
			Ok:     false,
			Errors: []string{"Не существующая комната"},
		}
		return
	}
	room.engine.Exit()

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
	if usr.Room != "" {
		conn.MessageChan <- common.Message{
			Type:   common.HUB_SUBSCRIBE,
			Ok:     false,
			Errors: []string{"Вы уже подписаны на другую комнату"},
		}
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
		if err := room.engine.AddParticipant(enmodels.Participant{
			Name: usr.Name,
			Role: usr.Role,
		}); err != nil {
			if enerr.KindIs(enerr.AlreadyStarted, err) {
				conn.MessageChan <- common.Message{
					Type:   common.HUB_SUBSCRIBE,
					Ok:     false,
					Errors: []string{"Комната уже запущена, в неё нельзя войти"},
				}
			}
			if enerr.KindIs(enerr.MaxPlayers, err) {
				conn.MessageChan <- common.Message{
					Type:   common.HUB_SUBSCRIBE,
					Ok:     false,
					Errors: []string{"В комнате заняты все места"},
				}
			}
			return
		}

		oldRoomID := usr.setRoom(data.Name)
		oldRoom, ok := h.Rooms.get(oldRoomID)
		if ok {
			if err := oldRoom.engine.RemoveParticipant(e.userId); err != nil {
				log.Error("Failed to delete player", sl.Err(err))
				return
			}
		} else {
			log.Error("Failed to get room")
		}
		h.Channels.Remove(oldRoomID, connID)
		h.Channels.Add(data.Name, connID)
		go h.SendMessageOverChannel(data.Name, common.Message{
			Type: common.ENGINE_INFO,
			Ok:   true,
			Body: room.engine.PreHook(),
		})
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
func UnSubscribe(h *Hub, e internalEventWrap) {
	type dataT struct {
		Type   string
		Target string
		Name   string
	}
	op := "polymers/UnSubscribe"
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

	log.Debug("Context for UnSubscribe Event", "user", usr.Name, "data", data)
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
				Type:   common.HUB_UNSUBSCRIBE,
				Ok:     false,
				Errors: []string{"Не существующая комната"},
			}
			return
		}
		if err := room.engine.RemoveParticipant(e.userId); err != nil {
			if enerr.KindIs(enerr.AlreadyStarted, err) {
				conn.MessageChan <- common.Message{
					Type:   common.HUB_UNSUBSCRIBE,
					Ok:     false,
					Errors: []string{"Комната уже запущена, из нее нельзя выйти"},
				}
			} else {
				log.Debug("Failed Remove player", sl.Err(err))
			}
			return
		}
		usr.Room = ""
		usr.channels = remove(usr.channels, data.Name)
		h.Channels.Remove(data.Name, connID)
	case "channel":
		usr.channels = remove(usr.channels, data.Name)
		h.Channels.Remove(data.Name, connID)
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
		Type: common.HUB_UNSUBSCRIBE,
		Ok:   true,
		Body: map[string]any{
			"Target": "room",
			"Name":   data.Name,
		},
	}

}
func EngineAction(h *Hub, e internalEventWrap) {
	room, ok := h.Rooms.get(e.room)
	if !ok {
		h.log.Error("Cannot find room for EngineEvent", "room", e.room)
		return
	}
	go room.engine.Input(enmodels.Action{
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
	if e.role < common.Judge_Role {
		log.Error("No permission")
		return
	}
	room, ok := h.Rooms.get(e.room)
	if !ok {
		log.Error("Failed getting room", "roomname", e.room)
		return
	}
	room.engine.Start()
}
