package polymers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_checkIfHasElements(t *testing.T) {
	// Создаем тестовые данные
	player := &Player{
		Bag: map[string]int{
			"wood":  5,
			"stone": 3,
			"gold":  0,
		},
	}

	tests := []struct {
		name           string
		inputElements  map[string]int
		expectedResult []string
		expectedError  bool
	}{
		{
			name: "valid_elements",
			inputElements: map[string]int{
				"wood":  4,
				"stone": 2,
			},
			expectedResult: []string{},
			expectedError:  false,
		},
		{
			name: "invalid_elements_missing",
			inputElements: map[string]int{
				"wood":  6, // Недостаточно
				"stone": 2,
			},
			expectedResult: []string{"wood"},
			expectedError:  true,
		},
		{
			name: "invalid_elements_not_in_bag",
			inputElements: map[string]int{
				"diamond": 1, // Элемента нет в сумке
			},
			expectedResult: []string{"diamond"},
			expectedError:  true,
		},
		{
			name:           "empty_input",
			inputElements:  map[string]int{},
			expectedResult: []string{},
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := player.checkIfHasElements(tt.inputElements)

			// Проверка результата
			assert.ElementsMatch(t, tt.expectedResult, result)

			// Проверка ошибки
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
