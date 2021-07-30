package service

import (
	"encoding/json"
	"io"

	"github.com/qbart/go-grpc/models"
	"go.uber.org/zap"
)

// PortsReader decodes json file with ports,
// then each decoded Port is sent via channel
//
type PortsReader struct {
	Logger *zap.Logger
}

func (pr *PortsReader) Stream(reader io.Reader) <-chan *models.Port {
	ch := make(chan *models.Port)

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
				var port models.Port
				err := decoder.Decode(&port)
				if err != nil {
					pr.Logger.Error("Can't decode Port")
					break
				}
				port.ID = id

				ch <- &port
			}
		}
	}(reader)

	return ch
}
