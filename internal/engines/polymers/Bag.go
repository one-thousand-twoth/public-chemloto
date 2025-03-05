package polymers

import (
	"encoding/json"
	"errors"
	"maps"
	"math/rand"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
)

var (
	ErrEmptyBag = errors.New("empty bag")
)

type Bag struct {
}

type GameBag struct {
	// Общий мешок key = название элемента, value = количество
	Elements map[string]int
	// Текущие элементы в мешке
	iterElements map[string]int
	// Вытащенные элементы
	PushedValues []string
	rnd          *rand.Rand
}

func (b GameBag) MarshalJSON() ([]byte, error) {
	bag := struct {
		Elements     map[string]int
		LastElements []string
	}{
		b.Elements,
		b.LastElements(),
	}
	return json.Marshal(bag)
}

// NewGameBag initializes GameBag with elements and generate random seed.
func NewGameBag(elements map[string]int) GameBag {
	keys := make([]string, 0, 12)
	for k, v := range elements {
		if v != 0 {
			for i := 0; i < v; i++ {
				keys = append(keys, k)
			}
		}
	}
	return GameBag{
		Elements:     elements,
		iterElements: maps.Clone(elements),
		PushedValues: make([]string, 0, len(keys)),
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (b *GameBag) getRandomElement() (string, error) {
	const op enerr.Op = "polymers/GameBag.getRandomElement"
	element, err := pickRandom(b.iterElements, b.rnd)
	if err != nil {
		return "", enerr.E(op, err)
	}
	b.PushedValues = append(b.PushedValues, element)
	return element, nil
}

// pickRandom возвращает случайный элемент из from
func pickRandom(from map[string]int, rand *rand.Rand) (string, error) {
	const op enerr.Op = "polymers/GameBag.pickRandom"
	// Вычисляем общее количество
	total := 0
	for _, count := range from {
		total += count
	}

	// Если ничего не осталось
	if total == 0 {
		return "", enerr.E(op, ErrEmptyBag)
	}

	randomNumber := rand.Intn(total)

	// Определяем, какой ключ выбран, и уменьшаем его количество
	cumulative := 0
	for k, count := range from {
		cumulative += count
		if randomNumber < cumulative {
			from[k]--
			return k, nil
		}
	}

	return "", enerr.E(op, "Ошибка при вытаскивании нового элемента из мешка", enerr.Internal)
}

func (b *GameBag) LastElements() []string {
	if len(b.PushedValues) < 6 {
		return b.PushedValues
	}
	return b.PushedValues[len(b.PushedValues)-6:]
}
