package main

import (
	"context"
	"log"
	"net"

	"github.com/qbart/go-grpc/pb"
	"google.golang.org/grpc"
)

type portDomainService struct {
	pb.UnimplementedPortDomainServiceServer
}

func (s *portDomainService) Upsert(ctx context.Context, port *pb.PortRequest) (*pb.PortResponse, error) {
	log.Println(port)
	return nil, nil
}

func newServer() *portDomainService {
	s := &portDomainService{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPortDomainServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
