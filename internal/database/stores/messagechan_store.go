package stores

import (
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
)

func NewStreamStore() *StreamStore {
	return &StreamStore{
		engines: map[string]entities.MessageStream{},
		mu:      sync.Mutex{},
	}
}

type StreamStore struct {
	engines map[string]entities.MessageStream
	mu      sync.Mutex
}

func (store *StreamStore) Add(name string, mss entities.MessageStream) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.engines[name] = mss
}

func (store *StreamStore) Get(name string) entities.MessageStream {
	store.mu.Lock()
	defer store.mu.Unlock()

	return store.engines[name]
}
