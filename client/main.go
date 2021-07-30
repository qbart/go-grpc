package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qbart/go-grpc/client/web"
	"github.com/qbart/go-grpc/pb"
	"github.com/qbart/go-grpc/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	grpcAddr := envOrDefault("PORTS_SERVICE_ADDR", "0.0.0.0:3001")
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		logger.Warn("Retrying connection to grpc in 2 secs...")
		time.Sleep(time.Second * 2)

		conn, err = grpc.Dial(grpcAddr, grpc.WithInsecure())
		if err != nil {
			logger.Fatal("Faild to connect to grpc", zap.Error(err))
		}
	}
	defer conn.Close()
	portDomainService := pb.NewPortDomainServiceClient(conn)

	portsSyncronizer := service.PortsSyncronizer{
		Logger:            logger,
		ConfigFileName:    envOrDefault("CONFIG_FILE", "ports.json"),
		PortDomainService: portDomainService,
	}

	// start syncronization in the background
	// so http server can start immediately
	go func() {
		err := portsSyncronizer.Sync()
		if err != nil {
			logger.Fatal("Port syncronization failed to start", zap.Error(err))
		}
	}()

	portAPI := web.PortAPIHandler{
		Logger:            logger,
		PortDomainService: portDomainService,
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/ports/{id:[A-Z]+}", portAPI).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:3000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info("Closing http server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-c

	logger.Info("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Failed to perform a shutdown", zap.Error(err))
	}
}

func envOrDefault(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}
	return value
}
