package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/qbart/go-grpc/pb"
)

type MemDB struct {
	ports map[string]*pb.Port
	mu    sync.RWMutex
}

func New() *MemDB {
	return &MemDB{
		ports: make(map[string]*pb.Port),
	}
}

func (db *MemDB) Upsert(ctx context.Context, port *pb.Port) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.ports[port.Id] = port

	return nil
}

func (db *MemDB) Get(ctx context.Context, id string) (*pb.Port, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if port, ok := db.ports[id]; ok {
		return port, nil
	}

	return nil, fmt.Errorf("Port with ID=%s does not exist", id)
}
