package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qbart/go-grpc/models"
	"github.com/qbart/go-grpc/storage"
)

func main() {
	ctx := context.Background()
	db := storage.New()

	go func() {
		f, err := os.Open("ports.json")
		if err != nil {
			log.Fatalf("Can't open file: %v", err)
		}
		defer f.Close()

		ports := NewPortsReader(f)
		for port := range ports {
			db.Upsert(port)
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
}

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

func PortHandler(db storage.DB) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		port, err := db.Get(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		err = json.NewEncoder(w).Encode(port)
		if err != nil {
			// TODO
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
