package polymers

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrEmptyBag    = errors.New("empty bag")
	ErrInternalErr = errors.New("internal error")
)

type Bag struct {
}

type GameBag struct {
	// Общий мешок key = название элемента, value = количество
	Elements map[string]int
	// Названия элементов в количестве, необходим для вычисления случайного элемента
	Values []string
	// Вытащенные элементы
	PushedValues []string
	rnd          *rand.Rand
}

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
		Values:       keys,
		PushedValues: make([]string, 0, len(keys)),
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (b *GameBag) getRandomElement() (string, error) {

	if len(b.Values) == 0 {
		return "", ErrEmptyBag
	}

	rand_index := b.rnd.Intn(len(b.Values))
	elem := b.Values[rand_index]
	item, ok := b.Elements[elem]
	if !ok {
		// log.Error("Failed to pick an element", slog.String("Element", elem))
		return "", ErrInternalErr
	}

	b.Elements[elem] = item - 1
	b.Values = removeByValue(b.Values, elem)
	b.PushedValues = append(b.PushedValues, elem)

	return elem, nil
}

func (b *GameBag) LastElements() []string {
	if len(b.PushedValues) < 5 {
		return b.PushedValues
	}
	return b.PushedValues[len(b.PushedValues)-5:]
}
