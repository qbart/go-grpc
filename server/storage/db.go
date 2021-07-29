package storage

import (
	"context"

	"github.com/qbart/go-grpc/pb"
)

type DB interface {
	Upsert(ctx context.Context, port *pb.Port) error
	Get(ctx context.Context, id string) (*pb.Port, error)
}
