package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qbart/go-grpc/models"
	"github.com/qbart/go-grpc/pb"
	"github.com/qbart/go-grpc/storage"
	"google.golang.org/grpc"
)

func main() {
	db := storage.New()

	conn, err := grpc.Dial("0.0.0.0:3001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to grpc %v", err)
	}
	defer conn.Close()
	client := pb.NewPortDomainServiceClient(conn)

	go func() {
		f, err := os.Open("ports.json")
		if err != nil {
			log.Fatalf("Can't open file: %v", err)
		}
		defer f.Close()

		ports := NewPortsReader(f)
		for port := range ports {
			_, err := client.Upsert(context.Background(), &pb.PortRequest{
				Id:          "port.ID" + port.ID,
				Name:        "port.Name",
				City:        "port.City",
				Country:     "port.Country",
				Alias:       []string{"aa"},
				Regions:     []string{"aa"},
				Coordinates: []float64{1, 2},
				Province:    "port.Province",
				Timezone:    "port.Timezone",
				Unlocs:      []string{"aa"},
				// Id:          port.ID,
				// Name:        port.Name,
				// City:        port.City,
				// Country:     port.Country,
				// Alias:       port.Alias,
				// Regions:     port.Regions,
				// Coordinates: port.Coordinates[:],
				// Province:    port.Province,
				// Timezone:    port.Timezone,
				// Unlocs:      port.Unlocs,
			})
			if err != nil {
				log.Fatalf("Upserting port failed: %v", err)
			}
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/ports/{id:[A-Z]+}", PortHandler(db)).Methods("GET")

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

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

func PortHandler(db storage.DB) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		port, err := db.Get(params["id"])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		if err = json.NewEncoder(w).Encode(port); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
	}
}

func NewPortsReader(r io.Reader) <-chan *models.Port {
	ch := make(chan *models.Port)

	go func(r io.Reader) {
		defer close(ch)

		dec := json.NewDecoder(r)

		for {
			t, err := dec.Token()
			if err == io.EOF {
				break
			}

			if id, ok := t.(string); ok {
				var port models.Port
				err = dec.Decode(&port)
				port.ID = id
				ch <- &port
			}
		}
	}(r)

	return ch
}
