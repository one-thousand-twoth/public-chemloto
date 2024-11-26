package hub

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/appvalidation"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	enmodels "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
	"github.com/anrew1002/Tournament-ChemLoto/web"
)

type Engine interface {
	// Получить текущее состояние, например при перезагрузке страницы
	PreHook() map[string]any
	// Обработать событие
	Input(enmodels.Action)
	Start()
	GetResults() [][]string
	AddPlayer(enmodels.Player) error
	RemovePlayer(name string) error
}

func (h *Hub) AddNewRoom(r Room) error {

	validate := appvalidation.Ins
	if err := validate.Struct(r); err != nil {
		// validateErr := err.(valid.ValidationErrors)
		return err
	}
	if !r.IsAuto {
		r.Time = 0
	}
	checks := parseEngineJson(h)
	r.engine = polymers.New(
		h.log.With(slog.String("room", r.Name)),
		polymers.PolymersEngineConfig{
			Elements:   r.Elements,
			Checks:     checks,
			TimerInt:   r.Time,
			MaxPlayers: r.MaxPlayers,
			Unicast: func(userID string, msg common.Message) {
				go func() {
					h.log.Debug("Unicast message")
					usr, ok := h.Users.Get(userID)
					if !ok {
						h.log.Error("failed to get user while Unicast message from engine")
						return
					}
					connID := usr.GetConnection()
					conn, ok := h.Connections.Get(connID)
					if !ok {
						h.log.Error("failed to get user connection while Unicast message from engine")
						return
					}
					conn.MessageChan <- msg
				}()
			},
			Broadcast: func(msg common.Message) {
				h.log.Debug("Broadcast message")
				go h.SendMessageOverChannel(r.Name, msg)
			},
		},
	)

	if err := h.Rooms.add(&r); err != nil {
		return err
	}
	// Set Init function for all room channels
	h.Channels.SetChannelFunc(r.Name, func(ch chan common.Message) {
		body := r.engine.PreHook()
		// h.log.Error("{Unsafe} sending engine state", "body", body)
		ch <- common.Message{
			Type: common.ENGINE_INFO,
			Ok:   true,
			Body: body,
		}
	})
	return nil
}

func parseEngineJson(h *Hub) polymers.Checks {
	jsonFile, err := web.Polymers.Open("polymers.json")
	if err != nil {
		h.log.Error("Cannot open polymers config", sl.Err(err))
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		h.log.Error("Error reading polymers", sl.Err(err))
	}
	var checks polymers.Checks
	err = json.Unmarshal(byteValue, &checks.Fields)
	if err != nil {
		h.log.Error("polymers.json have errors", sl.Err(err))
	}
	return checks
}

type Room struct {
	Name       string         `json:"name" validate:"required,min=1,safeinput"`
	MaxPlayers int            `json:"maxPlayers" validate:"required,gt=1,lt=100"`
	Elements   map[string]int `json:"elementCounts" validate:"required"`
	Time       int            `validate:"excluded_if=isAuto false,gte=0"`
	IsAuto     bool           `json:"isAuto"`
	engine     Engine         `json:"-"`
}

type roomsState struct {
	// map key is a room name
	state map[string]*Room
	mutex sync.RWMutex
}

func (rs *roomsState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
}

func (rs *roomsState) get(id string) (r *Room, ok bool) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	r, ok = rs.state[id]
	return
}
func (rs *roomsState) add(room *Room) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	_, ok := rs.state[room.Name]
	if ok {
		return errors.New("already exist room")
	} else {
		rs.state[room.Name] = room
	}

	return nil
}
