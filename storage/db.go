package storage

import (
	"github.com/qbart/go-grpc/models"
)

type DB interface {
	Upsert(port *models.Port) error
	Get(id string) (*models.Port, error)
}
