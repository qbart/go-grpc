package serializers

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/qbart/go-grpc/models"
	"github.com/qbart/go-grpc/pb"
)

func TestPortSerialization(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#SerializePort", func() {
		g.It("Serializes Port model to protobuf", func() {
			port := models.Port{
				ID:          "1",
				Name:        "Name",
				City:        "City",
				Country:     "Country",
				Alias:       []string{"alias"},
				Regions:     []string{"regions"},
				Coordinates: []float64{1.0, 2.0},
				Province:    "Province",
				Unlocs:      []string{"Unlocs"},
				Code:        "Code",
			}

			pbPort := SerializePort(&port)

			g.Assert(pbPort.Id).Eql("1")
			g.Assert(pbPort.Name).Eql("Name")
			g.Assert(pbPort.City).Eql("City")
			g.Assert(pbPort.Country).Eql("Country")
			g.Assert(pbPort.Alias).Eql([]string{"alias"})
			g.Assert(pbPort.Regions).Eql([]string{"regions"})
			g.Assert(pbPort.Coordinates).Eql([]float64{1.0, 2.0})
			g.Assert(pbPort.Province).Eql("Province")
			g.Assert(pbPort.Unlocs).Eql([]string{"Unlocs"})
			g.Assert(pbPort.Code).Eql("Code")
		})
	})

	g.Describe("#DeserializePort", func() {
		g.It("Deserializes proto Port to model", func() {
			pbPort := pb.Port{
				Id:          "1",
				Name:        "Name",
				City:        "City",
				Country:     "Country",
				Alias:       []string{"alias"},
				Regions:     []string{"regions"},
				Coordinates: []float64{1.0, 2.0},
				Province:    "Province",
				Unlocs:      []string{"Unlocs"},
				Code:        "Code",
			}

			port := DeserializePort(&pbPort)

			g.Assert(port.ID).Eql("1")
			g.Assert(port.Name).Eql("Name")
			g.Assert(port.City).Eql("City")
			g.Assert(port.Country).Eql("Country")
			g.Assert(port.Alias).Eql([]string{"alias"})
			g.Assert(port.Regions).Eql([]string{"regions"})
			g.Assert(port.Coordinates).Eql([]float64{1.0, 2.0})
			g.Assert(port.Province).Eql("Province")
			g.Assert(port.Unlocs).Eql([]string{"Unlocs"})
			g.Assert(port.Code).Eql("Code")
		})
	})
}
