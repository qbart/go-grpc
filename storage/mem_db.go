package storage

import (
	"fmt"
	"sync"

	"github.com/qbart/go-grpc/models"
)

type MemDB struct {
	ports map[string]*models.Port
	mu    sync.RWMutex
}

func New() *MemDB {
	return &MemDB{
		ports: make(map[string]*models.Port),
	}
}

func (db *MemDB) Upsert(port *models.Port) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.ports[port.ID] = port

	return nil
}

func (db *MemDB) Get(id string) (*models.Port, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if port, ok := db.ports[id]; ok {
		return port, nil
	}

	return nil, fmt.Errorf("Port with ID=%s does not exist", id)
}
