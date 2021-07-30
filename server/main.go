package main

import (
	"context"
	"net"

	"github.com/qbart/go-grpc/pb"
	"github.com/qbart/go-grpc/serializers"
	"github.com/qbart/go-grpc/server/storage"
	"github.com/qbart/go-grpc/server/storage/memory"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type portDomainService struct {
	pb.UnimplementedPortDomainServiceServer

	db     storage.DB
	logger *zap.Logger
}

func (s *portDomainService) Upsert(ctx context.Context, port *pb.Port) (*pb.UpsertResponse, error) {
	s.logger.Debug("Upsert", zap.String("PortID", port.Id))

	err := s.db.Upsert(ctx, serializers.DeserializePort(port))
	return &pb.UpsertResponse{}, err
}

func (s *portDomainService) Get(ctx context.Context, portId *pb.PortId) (*pb.Port, error) {
	s.logger.Debug("Get", zap.String("PortID", portId.Id))

	port, err := s.db.Get(ctx, portId.Id)
	if err != nil {
		return nil, err
	}
	return serializers.SerializePort(port), nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db := memory.New()

	listener, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPortDomainServiceServer(grpcServer, &portDomainService{db: db, logger: logger})
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal("GRPC server error", zap.Error(err))
	}
}
