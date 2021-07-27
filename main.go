package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Coordinates = [2]float64

type Port struct {
	ID          string
	Name        string      `json:"name"`
	City        string      `json:"city"`
	Alias       []string    `json:"alias"`
	Regions     []string    `json:"regions"`
	Coordinates Coordinates `json:"coordinates"`
	Province    string      `json:"province"`
	Timezone    string      `json:"timezone"`
	Unlocs      []string    `json:"unlocs"`
}

func main() {
	f, err := os.Open("ports.json")
	if err != nil {
		log.Fatalf("Can't open file: %v", err)
	}
	defer f.Close()

	ports := NewPortsReader(f)
	for port := range ports {
		fmt.Println(port)
	}
}

func NewPortsReader(r io.Reader) <-chan *Port {
	ch := make(chan *Port)

	go func(r io.Reader) {
		defer close(ch)

		dec := json.NewDecoder(r)

		for {
			t, err := dec.Token()
			if err == io.EOF {
				break
			}

			if id, ok := t.(string); ok {
				var port Port
				err = dec.Decode(&port)
				port.ID = id
				ch <- &port
			}
		}
	}(r)

	return ch
}
