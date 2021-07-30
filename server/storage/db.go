package storage

import (
	"context"

	"github.com/qbart/go-grpc/models"
)

type DB interface {
	Upsert(ctx context.Context, port *models.Port) error
	Get(ctx context.Context, id string) (*models.Port, error)
}
