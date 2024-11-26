package polymers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models/enerr"
)

var (
	ErrEmptyBag = errors.New("empty bag")
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
		Values:       keys,
		PushedValues: make([]string, 0, len(keys)),
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (b *GameBag) getRandomElement() (string, error) {
	const op enerr.Op = "polymers/GameBag.getRandomElement"
	if len(b.Values) == 0 {

		return "", enerr.E(op, ErrEmptyBag)
	}

	randIndex := b.rnd.Intn(len(b.Values))
	elem := b.Values[randIndex]
	item, ok := b.Elements[elem]
	if !ok {
		// log.Error("Failed to pick an element", slog.String("Element", elem))
		return "", enerr.E(op, "Ошибка при вытаскивании нового элемента из мешка", enerr.Internal)
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
