package polymers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockEngine creates a mock PolymersEngine for testing
func MockEngine() *PolymersEngine {
	return &PolymersEngine{
		maxPlayers: 4,
		// log:        &mockLogger{},
	}
}

// type mockLogger struct{}

// func (m *mockLogger) Debug(msg string, args ...any) {}
// func (m *mockLogger) Error(msg string, args ...any) {}
// func (m *mockLogger) Info(msg string, args ...any)  {}
// func (m *mockLogger) Warn(msg string, args ...any)  {}

// func (m *mockLogger) With(args ...any) common.Logger { return m }

func TestSimpleState(t *testing.T) {
	t.Run("Add Handler", func(t *testing.T) {
		state := NewState()
		handler := func(a models.Action) (stateInt, error) { return NO_TRANSITION, nil }

		state = state.Add("TestAction", handler, false)

		assert.Contains(t, state.handlers, "TestAction")
		assert.False(t, state.secure["TestAction"])
	})

	t.Run("Secure Handler", func(t *testing.T) {
		state := NewState()
		handler := func(a models.Action) (stateInt, error) { return NO_TRANSITION, nil }

		state = state.Add("AdminAction", handler, true)

		assert.Contains(t, state.handlers, "AdminAction")
		assert.True(t, state.secure["AdminAction"])
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		state := NewState()

		jsonData, err := state.MarshalJSON()

		require.NoError(t, err)
		assert.Equal(t, []byte("null"), jsonData)
	})
}

func TestObtainState(t *testing.T) {
	t.Run("Create ObtainState", func(t *testing.T) {
		eng := MockEngine()
		obtainState := eng.NewObtainState(10*time.Second, false)

		assert.NotNil(t, obtainState)
		assert.Contains(t, obtainState.handlers, "GetElement")
		assert.Contains(t, obtainState.handlers, "RaiseHand")
	})
}

func TestTradeState(t *testing.T) {
	t.Run("Create TradeState", func(t *testing.T) {
		eng := MockEngine()
		tradeState := eng.NewTradeState(10 * time.Second)

		assert.NotNil(t, tradeState)
		assert.NotNil(t, tradeState.StockExchange)
		assert.Contains(t, tradeState.handlers, "TradeOffer")
		assert.Contains(t, tradeState.handlers, "TradeRequest")
		assert.Contains(t, tradeState.handlers, "TradeAck")
	})

	t.Run("StockExchange Operations", func(t *testing.T) {
		se := &StockExchange{
			StockList: make([]*Stock, 0),
			TradeLog:  make([]*TradeLog, 0),
		}

		// Test AddStock
		stock := &Stock{
			ID:       "test-stock",
			Requests: make(map[string]*StockRequest),
		}
		se.AddStock("test-stock", stock)
		assert.Len(t, se.StockList, 1)

		// Test StockByID
		foundStock, err := se.StockByID("test-stock")
		require.NoError(t, err)
		assert.Equal(t, stock, foundStock)

		// Test SetRequest
		request := &StockRequest{
			ID:     "test-request",
			Player: "test-player",
			Accept: true,
		}
		err = se.SetRequest("test-stock", request)
		require.NoError(t, err)
		assert.Contains(t, foundStock.Requests, "test-player")
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		eng := MockEngine()
		tradeState := eng.NewTradeState(10 * time.Second)

		jsonData, err := tradeState.MarshalJSON()

		require.NoError(t, err)

		var result struct {
			StockExchange StockExchange
			Timer         int
		}
		err = json.Unmarshal(jsonData, &result)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.Timer, 0)
	})
}

// // Helper function to parse action data
// func dataFromAction[T any](action models.Action) (T, error) {
// 	var data T
// 	// Implementation would depend on how dataFromAction is actually implemented in the original code
// 	// This is a mock implementation for testing
// 	return data, nil
// }
