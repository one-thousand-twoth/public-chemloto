package hub

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/usecase"
	"github.com/mitchellh/mapstructure"
)

// NOTE: лучше передавать канал для сообщений пользователю
type HandlerFunc2 func(internalEventWrap) error

type WebsocketHandlers struct {
	eventHandlers map[string]HandlerFunc2

	usecases *usecase.Usecases

	log *slog.Logger
}

func NewWebsocketHandlers(
	usecases *usecase.Usecases,
	log *slog.Logger,
) *WebsocketHandlers {

	wh := &WebsocketHandlers{
		eventHandlers: map[string]HandlerFunc2{},
		log:           log,
		usecases:      usecases,
	}
	wh.SetupHandlers()
	return wh

}
func (wh *WebsocketHandlers) Handle(e internalEventWrap) error {
	handler, ok := wh.eventHandlers[e.msgType.String()]
	if !ok {
		// h.log.Error("error find event handler ", "message type", event.msgType.String())
		return enerr.E("Cannot find handler")
	}

	go func() {
		err := handler(e)
		if err != nil {
			wh.log.Error("Unhandled error", "err", err)
		}
	}()

	return nil
}

// func (f HandlerFunc) HandleEvent(h *Hub, msg internalEventWrap) {
// 	f(h, msg)
// }

// use handler f for messages t
func (h *WebsocketHandlers) UseHandler(t common.MessageType, f HandlerFunc2) {
	h.eventHandlers[t.String()] = f
}

func (h *WebsocketHandlers) SetupHandlers() {
	h.UseHandler(common.HUB_SUBSCRIBE, h.SubscribeHandler2)
	h.UseHandler(common.HUB_UNSUBSCRIBE, h.UnSubscribeHandler2)
	h.UseHandler(common.ENGINE_ACTION, h.EngineAction2)
	h.UseHandler(common.HUB_STARTGAME, h.StartGame2)
	h.UseHandler(common.HUB_EXITGAME, h.StopGame)

}

type SubscribeRequest struct {
	Type   string
	Target string
	Name   string
}

// Subscribe обрабатывает логику добавления user`a в подписчики канала
func (h *WebsocketHandlers) SubscribeHandler2(e internalEventWrap) error {

	const op enerr.Op = "Subscribe handler"
	log := h.log.With(enerr.OpAttr(op))

	log.Debug("Start Handle Event", "usr", e.userId, "data", fmt.Sprintf("%v", e.msg))

	var data SubscribeRequest
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		return enerr.E(op, fmt.Sprintf("failed to decode event body: %s", err.Error()))
	}

	switch data.Target {
	case "room":
		err := h.usecases.SubscribeToRoom(context.TODO(), data.Name, e.userId)
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
	// fun, ok := h.channelRepo.GetGroupByID()(data.Name)
	// if ok {
	// 	fun(e.MessageChannel)
	// }
	return nil
}

// Subscribe обрабатывает логику добавления user`a в подписчики канала
func (h *WebsocketHandlers) UnSubscribeHandler2(e internalEventWrap) error {
	type dataT struct {
		Type   string
		Target string
		Name   string
	}
	const op enerr.Op = "UnSubscribe handler"
	log := h.log.With("op", op)

	log.Debug("Start Handle Event", "usr", e.userId, "data", fmt.Sprintf("%v", e.msg))

	var data dataT
	if err := mapstructure.Decode(e.msg, &data); err != nil {
		return enerr.E(op, fmt.Sprintf("failed to decode event body: %s", err.Error()))
	}

	switch data.Target {
	case "room":
		err := h.usecases.UnsubscribeFromRoom(context.TODO(), data.Name, e.userId)
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
		Type: common.HUB_UNSUBSCRIBE,
		Ok:   true,
		Body: map[string]any{
			"Target": "room",
			"Name":   data.Name,
		},
	}

	return nil
}

func (h *WebsocketHandlers) EngineAction2(e internalEventWrap) error {
	const op enerr.Op = "EngineAction2 handler"

	err := h.usecases.RouteActionToUserRoom(context.TODO(), e.userId, e.msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *WebsocketHandlers) StartGame2(e internalEventWrap) error {
	type dataT struct {
		Type string
		Name string
	}
	op := "StartGame handler"
	log := h.log.With("op", op)
	log.Debug("Start Handle Event", "usr", e.userId, "data", fmt.Sprintf("%v", e.msg))

	h.usecases.StartGame(e.userId)

	return nil
}

func (h *WebsocketHandlers) StopGame(e internalEventWrap) error {
	type dataT struct {
		Type string
		Name string
	}
	op := "Hub/StopGame"
	log := h.log.With("op", op)
	var data dataT
	if err := mapstructure.Decode(e.msg, &data); err != nil {

		return enerr.E(op, fmt.Sprintf("failed to decode event body: %s", err.Error()))
	}
	log.Debug("Start Handle Event StopGame", "usr", e.userId, "room", data.Name, "data", fmt.Sprintf("%v", e.msg))

	err := h.usecases.StopGame(context.TODO(), data.Name, e.userId)
	if err != nil {
		return enerr.E(op, err)
	}
	return nil
}
