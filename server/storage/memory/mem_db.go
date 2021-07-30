package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qbart/go-grpc/models"
	"github.com/qbart/go-grpc/server/storage"
)

// memDB is a dummy in-memory database implementation.
type memDB struct {
	ports map[string]*models.Port
	mu    sync.RWMutex
}

func New() storage.DB {
	return &memDB{
		ports: make(map[string]*models.Port),
	}
}

func (db *memDB) Upsert(ctx context.Context, port *models.Port) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.ports[port.ID] = port

	return nil
}

func (db *memDB) Get(ctx context.Context, id string) (*models.Port, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if port, ok := db.ports[id]; ok {
		return port, nil
	}

	return nil, fmt.Errorf("Port with ID=%s does not exist", id)
}
