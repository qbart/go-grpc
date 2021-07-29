package main

import (
	"context"
	"log"
	"net"

	"github.com/qbart/go-grpc/pb"
	"github.com/qbart/go-grpc/storage"
	"google.golang.org/grpc"
)

type portDomainService struct {
	pb.UnimplementedPortDomainServiceServer

	db storage.DB
}

func (s *portDomainService) Upsert(ctx context.Context, port *pb.Port) (*pb.UpsertResponse, error) {
	log.Println("Upsert", port)
	err := s.db.Upsert(ctx, port)
	return &pb.UpsertResponse{}, err
}

func (s *portDomainService) Get(ctx context.Context, portId *pb.PortId) (*pb.Port, error) {
	log.Println("Get", portId)
	port, err := s.db.Get(ctx, portId.Id)
	return port, err
}

func main() {
	db := storage.New()

	lis, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPortDomainServiceServer(grpcServer, &portDomainService{db: db})
	grpcServer.Serve(lis)
}
