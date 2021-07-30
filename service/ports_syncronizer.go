package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/qbart/go-grpc/pb"
	"github.com/qbart/go-grpc/serializers"
	"go.uber.org/zap"
)

type PortsSyncronizer struct {
	Logger            *zap.Logger
	ConfigFileName    string
	PortDomainService pb.PortDomainServiceClient
}

func (ps *PortsSyncronizer) Sync() error {
	portsFile, err := os.Open(ps.ConfigFileName)
	if err != nil {
		return fmt.Errorf("Can't open the file %v", err)
	}
	defer func() {
		err := portsFile.Close()
		if err != nil {
			ps.Logger.Error("Can't close the file", zap.Error(err))
		}
	}()

	portsReader := PortsReader{Logger: ps.Logger}
	for port := range portsReader.Stream(portsFile) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err := ps.PortDomainService.Upsert(ctx, serializers.SerializePort(port))
		if err != nil {
			ps.Logger.Error("Upserting port failed", zap.Error(err), zap.String("PortID", port.ID))
		}
	}

	return nil
}
