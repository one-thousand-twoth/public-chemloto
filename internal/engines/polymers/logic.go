package polymers

import "errors"

var (
	ErrEmptyBag    = errors.New("empty bag")
	ErrInternalErr = errors.New("internal error")
)

func (engine *PolymersEngine) getRandomElement() (string, error) {
	keys := engine.lastElementsKeys

	if len(keys) == 0 {
		return "nil", ErrEmptyBag
	}

	rand_index := engine.rnd.Intn(len(keys))
	elem := keys[rand_index]
	item, ok := engine.Elements[elem]
	if !ok {
		// log.Error("Failed to pick an element", slog.String("Element", elem))
		return "Error", ErrInternalErr
	}

	engine.Elements[elem] = item - 1
	engine.lastElementsKeys = removeByValue(engine.lastElementsKeys, elem)
	engine.pushedElements = append(engine.pushedElements, elem)

	return elem, nil
}

func removeByValue(slice []string, value string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}
