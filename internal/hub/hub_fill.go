package hub

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
	roomName := "Тест Рук"
	h.AddNewRoom(
		CreateRoomRequest{
			Name:       roomName,
			MaxPlayers: 10,
			Time:       10,
			IsAuto:     true,
			Elements:   elements,
		},
	)

}
