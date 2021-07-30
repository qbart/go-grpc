package serializers

import (
	"github.com/qbart/go-grpc/models"
	"github.com/qbart/go-grpc/pb"
)

func SerializePort(p *models.Port) *pb.Port {
	alias := make([]string, len(p.Alias))
	regions := make([]string, len(p.Regions))
	coordinates := make([]float64, len(p.Coordinates))
	unlocs := make([]string, len(p.Unlocs))

	copy(alias, p.Alias)
	copy(regions, p.Regions)
	copy(coordinates, p.Coordinates)
	copy(unlocs, p.Unlocs)

	return &pb.Port{
		Id:          p.ID,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       alias,
		Regions:     regions,
		Coordinates: coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      unlocs,
		Code:        p.Code,
	}
}

func DeserializePort(p *pb.Port) *models.Port {
	alias := make([]string, len(p.Alias))
	regions := make([]string, len(p.Regions))
	coordinates := make([]float64, len(p.Coordinates))
	unlocs := make([]string, len(p.Unlocs))

	copy(alias, p.Alias)
	copy(regions, p.Regions)
	copy(coordinates, p.Coordinates)
	copy(unlocs, p.Unlocs)

	return &models.Port{
		ID:          p.Id,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       alias,
		Regions:     regions,
		Coordinates: coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      unlocs,
		Code:        p.Code,
	}
}
