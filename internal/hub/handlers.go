package hub

import (
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
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
	room.Engine.Exit()

}

type SubscribeRequest struct {
	Type   string
	Target string
	Name   string
}

// Subscribe обрабатывает логику добавления user`a в подписчики канала
func SubscribeHandler(h *Hub, e internalEventWrap) error {

	const op enerr.Op = "Subscribe handler"
	log := h.log.With("op", op)

	log.Debug("Start Handle Event", "usr", e.userId, "room", e.room, "data", fmt.Sprintf("%v", e.msg))

	var data SubscribeRequest
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		return enerr.E(op, fmt.Sprintf("failed to decode event body: %s", err.Error()))
	}

	usr, ok := h.Users.Get(e.userId)
	if !ok {
		return enerr.E(op, "Error getting user")
	}
	usr.mutex.Lock()
	defer usr.mutex.Unlock()

	log.Debug("Context for Subscribe Event", "user", usr.Name, "data", data)

	connID := usr.conn
	conn, ok := h.Connections.Get(connID)
	if !ok {
		return enerr.E(op, "Failed getting connection of user")
	}

	switch data.Target {
	case "room":
		err := subscribeToRoom(h, data, usr, log, connID)
		if err != nil {
			return err
		}
	case "channel":
		subscribeToChannel(usr, data, h, connID)
	default:
		return enerr.E(op, fmt.Sprintf("unknown target %s", data.Target))
	}

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
		fun(conn.MessageChan)
	}
	return nil
}

func subscribeToChannel(usr *User, data SubscribeRequest, h *Hub, connID string) {
	usr.channels = append(usr.channels, data.Name)
	h.Channels.Add(data.Name, connID)
}

func subscribeToRoom(
	h *Hub,
	data SubscribeRequest,
	// conn *SockConnection,
	usr *User, // TODO: есть опасность data race (!), блокировка происходит в вызывающей функции
	log *slog.Logger,
	connID string,
) error {
	const op enerr.Op = "hub.handlers/subscribeToRoom"

	if data.Target == "" || data.Name == "" {
		return enerr.E(op, "empty field", enerr.Validation)
	}

	room, ok := h.Rooms.get(data.Name)
	if !ok {
		return enerr.E(op, "room doesnt exist", enerr.NotExist)
	}

	err := room.Engine.AddParticipant(enmodels.Participant{
		Name: usr.Name,
		Role: usr.Role,
	})
	if err != nil {
		return enerr.E(op, err)
	}

	oldRoomID := usr.setRoom(data.Name)

	h.Channels.Remove(oldRoomID, connID)
	h.Channels.Add(data.Name, connID)
	return nil
}

func Subscribe(h *Hub, e internalEventWrap) {
	err := SubscribeHandler(h, e)
	if err != nil {
		h.log.Error("wip: error not propely handled", sl.Err(err))
	}
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
		if err := room.Engine.RemoveParticipant(e.userId); err != nil {
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
	go room.Engine.Input(enmodels.Action{
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
	room.Engine.Start()
}
