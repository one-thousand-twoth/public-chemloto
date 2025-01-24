package polymers

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/samber/lo"
)

// Hand defines struct for RaisedHand
type Hand struct {
	Player    *Player
	Field     string
	Name      string
	Structure map[string]int
	Checked   bool
}

// Game field with Scoring counter
type Field struct {
	Score int
}

// decrementScore substract by one field.Score and returns old value.
//
// Use it to assign Score to the Player who compose the Field correctly
func (f *Field) decrementScore() int {
	o := f.Score
	f.Score -= 1
	if f.Score < 0 {
		f.Score = 0
	}
	return o
}

// a Player represents Participant with Player roles with engine specific data.
type Player struct {
	models.Participant
	RaisedHand      bool
	Bag             map[string]int
	Score           int // Game score only for Players
	CompletedFields []string
}

// setScore will decrease Participant`s setScore with min value = -2
func (p *Player) setScore(score int) {
	p.Score += score
	if p.Score < -2 {
		p.Score = -2
	}
}
func (p *Player) raiseHand() error {
	return nil
}

// checkIfHasElements проверяет есть ли у игрока достаточно elements в его сумке
func (p *Player) checkIfHasElements(elements map[string]int) ([]string, error) {
	invalid := lo.PickBy(elements, func(k string, v int) bool {
		return v > p.Bag[k]
	})
	if len(invalid) > 0 {
		return lo.Keys(invalid), enerr.E("Указаны элементы которых нет в таком количестве", enerr.GameLogic)
	}
	return lo.Keys(invalid), nil
}
