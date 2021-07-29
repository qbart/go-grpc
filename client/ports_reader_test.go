package client

import (
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/qbart/go-grpc/pb"
)

func TestPortsReader(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#Stream", func() {
		g.It("Decodes stream correctly", func() {
			reader := NewPortsReader(strings.NewReader(`
{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": ["Abu"],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}
	`))

			result := make([]*pb.Port, 0)
			for port := range reader.Stream() {
				result = append(result, port)
			}

			g.Assert(len(result)).Eql(2)

			ajman := result[0]
			g.Assert(ajman.Id).Eql("AEAJM")
			g.Assert(ajman.Name).Eql("Ajman")
			g.Assert(ajman.City).Eql("Ajman")
			g.Assert(ajman.Country).Eql("United Arab Emirates")
			g.Assert(ajman.Alias).Eql([]string{})
			g.Assert(ajman.Regions).Eql([]string{})
			g.Assert(ajman.Coordinates).Eql([]float64{
				55.5136433,
				25.4052165,
			})
			g.Assert(ajman.Province).Eql("Ajman")
			g.Assert(ajman.Timezone).Eql("Asia/Dubai")
			g.Assert(ajman.Unlocs).Eql([]string{"AEAJM"})
			g.Assert(ajman.Code).Eql("52000")

			abu := result[1]
			g.Assert(abu.Id).Eql("AEAUH")
			g.Assert(abu.Name).Eql("Abu Dhabi")
			g.Assert(abu.City).Eql("Abu Dhabi")
			g.Assert(abu.Country).Eql("United Arab Emirates")
			g.Assert(abu.Alias).Eql([]string{"Abu"})
			g.Assert(abu.Regions).Eql([]string{})
			g.Assert(abu.Coordinates).Eql([]float64{
				54.37,
				24.47,
			})
			g.Assert(abu.Province).Eql("Abu Z¸aby [Abu Dhabi]")
			g.Assert(abu.Timezone).Eql("Asia/Dubai")
			g.Assert(abu.Unlocs).Eql([]string{"AEAUH"})
			g.Assert(abu.Code).Eql("52001")
		})
	})
}
