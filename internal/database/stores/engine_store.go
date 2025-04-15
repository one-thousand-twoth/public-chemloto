package stores

import (
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
)

func NewEngineStore() *EnginesStore {
	return &EnginesStore{
		engines: map[string]models.Engine{},
		mu:      sync.Mutex{},
	}
}

type EnginesStore struct {
	engines map[string]models.Engine
	mu      sync.Mutex
}

func (store *EnginesStore) Add(name string, engine models.Engine) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.engines[name] = engine
}

func (store *EnginesStore) Get(name string) models.Engine {
	store.mu.Lock()
	defer store.mu.Unlock()

	return store.engines[name]
}
