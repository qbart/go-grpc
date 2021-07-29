package client

import (
	"encoding/json"
	"io"
	"log"

	"github.com/qbart/go-grpc/pb"
)

type PortsReader struct {
	reader io.Reader
}

func NewPortsReader(reader io.Reader) *PortsReader {
	return &PortsReader{reader: reader}
}

func (pr *PortsReader) Stream() <-chan *pb.Port {
	ch := make(chan *pb.Port)

	go func(r io.Reader) {
		defer close(ch)

		decoder := json.NewDecoder(r)

		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}

			// everytime we encounter non-delimter we should have a string key that points to the object
			if id, ok := token.(string); ok {
				var port pb.Port
				err := decoder.Decode(&port)
				if err != nil {
					log.Println("Error decoding port")
					break
				}
				port.Id = id

				ch <- &port
			}
		}
	}(pr.reader)

	return ch
}
