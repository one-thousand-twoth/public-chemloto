// hub отвечает за контроль над вебсокет клиентами
package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/models"
	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
	mapset "github.com/deckarep/golang-set/v2"
)

type Hub struct {
	log *slog.Logger

	Rooms *roomsState
	Users *usersState

	Connections *connectionsState
	Channels    *channelsState

	eventHandlers map[string]HandlerFunc
	eventChan     chan internalEventWrap

	// registerChan   chan Instruction
	// unregisterChan chan Instruction
	// broadcastChan  chan Message
	// castChan       chan CastMessage
	// storage        *sqlite.Storage
	// rooms          map[string]*Room
	// roomMutex      sync.RWMutex
	// freeUsers      map[string]*Client
	// usersMutex     sync.RWMutex
}

func NewHub(storage *sqlite.Storage, log *slog.Logger) *Hub {
	return &Hub{
		log:           log,
		Rooms:         &roomsState{state: make(map[string]*room)},
		Users:         &usersState{state: make(map[string]*User)},
		Connections:   &connectionsState{state: make(map[string]*SockConnection)},
		Channels:      &channelsState{state: make(map[string]mapset.Set[string])},
		eventHandlers: make(map[string]HandlerFunc),
		eventChan:     make(chan internalEventWrap),
	}
}

func (h *Hub) Run() {
	go func() {
	eventLoop:
		for event := range h.eventChan {
			handler, ok := h.eventHandlers[event.msgType.String()]
			if !ok {
				h.log.Error("error find event handler ", "message type", event.msgType.String())
				continue eventLoop
			}
			handler.HandleEvent(h, event)
		}
	}()
}
func (h *Hub) SendEventToHub(e internalEventWrap) {
	h.eventChan <- e
}

func (h *Hub) SendMessageOverChannel(channel string, message models.Message) error {
	op := "SendMessageOverChannel"
	log := h.log.With("op", op)
	connections, ok := h.Channels.Get(channel)
	if !ok {
		return fmt.Errorf("no channel with %s name", channel)
	}
	log.Debug(fmt.Sprintf("WS Message: %+v to %s", message, channel))
	for _, connectionName := range connections {
		conn, ok := h.Connections.Get(connectionName)
		if ok {
			conn.MessageChan <- message
		}
	}
	return nil
}

func GetMessageType(msg []byte) (map[string]interface{}, MessageType, error) {

	var payload map[string]interface{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		return nil, UNDEFINED, errors.New("fail to unmarshal struct")
	}
	payloadType, ok := payload["Type"]
	if !ok {
		return nil, UNDEFINED, errors.New("поле Type не указано")
	}
	msgType, ok := MapEnumStringToMessageType[payloadType.(string)]
	if !ok {
		return nil, UNDEFINED, errors.New("неизвестный тип сообщения")
	}
	// msgTyped, err := NewMessage(payloadType, payload)
	return payload, msgType, err
}
