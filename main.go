package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qbart/go-grpc/client"
	"github.com/qbart/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:3001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to grpc %v", err)
	}
	defer conn.Close()
	portDomainService := pb.NewPortDomainServiceClient(conn)

	go func() {
		portsFile, err := os.Open("ports.json")
		if err != nil {
			log.Fatalf("Can't open file: %v", err)
		}
		defer portsFile.Close()

		portsReader := client.NewPortsReader(portsFile)
		for port := range portsReader.Stream() {
			_, err := portDomainService.Upsert(context.Background(), port)
			if err != nil {
				log.Fatalf("Upserting port %v failed: %v", port, err)
			}
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/ports/{id:[A-Z]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		port, err := portDomainService.Get(context.Background(), &pb.PortId{Id: params["id"]})

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err = json.NewEncoder(w).Encode(port); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:3000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-c

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
}
