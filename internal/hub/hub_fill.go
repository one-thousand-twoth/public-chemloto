package hub

import (
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
)

func (h *Hub) FillRooms() {
	// elements := map[string]int{
	// 	"TRADE": 4,
	// 	"O":     28,
	// 	"N":     16,
	// 	"H":     52,
	// 	"Cl":    16,
	// 	"CH3":   28,
	// 	"CH2":   24,
	// 	"CH":    24,
	// 	"C6H5":  16,
	// 	"C6H4":  16,
	// 	"C":     40,
	// }
	elements := map[string]int{
		"TRADE": 16,
		"O":     0,
		"N":     0,
		"H":     0,
		"Cl":    16,
		"CH3":   0,
		"CH2":   0,
		"CH":    0,
		"C6H5":  0,
		"C6H4":  16,
		"C":     0,
	}
	checks := parseEngineJson(h)
	roomName := "Тест Рук"
	h.AddNewRoom(
		Room{
			Name:       roomName,
			MaxPlayers: 10,
			Time:       10,
			IsAuto:     true,
			Elements:   elements,
			Engine: polymers.New(
				h.log.With(slog.String("room", roomName)),
				polymers.PolymersEngineConfig{
					Elements:   elements,
					Checks:     checks,
					TimerInt:   10,
					MaxPlayers: 10,
					Unicast: func(userID string, msg common.Message) {
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
					},
					Broadcast: func(msg common.Message) {
						h.log.Debug("Broadcast message")
						h.SendMessageOverChannel(roomName, msg)
					},
				},
			),
		},
	)

}
