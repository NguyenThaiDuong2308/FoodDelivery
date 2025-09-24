package repository

import (
	"context"
	"log"
	"shipper-assignment-service/internal/models"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type geoOperations struct {
	client *redis.Client
	geoKey string
}

type GeoOperations interface {
	FindNearestShipper(ctx context.Context, coord *models.Coordinate, radiusKM float64) (int, error)
}

func NewGeoOperations(client *redis.Client, geoKey string) GeoOperations {
	return &geoOperations{
		client: client,
		geoKey: geoKey,
	}
}

func (g *geoOperations) FindNearestShipper(ctx context.Context, coord *models.Coordinate, radiusKM float64) (int, error) {
	log.Println("coordinate:", coord)
	log.Println("GEOKEY:", g.geoKey)
	log.Println("radius:", radiusKM)
	cmd := g.client.GeoRadius(ctx, g.geoKey, coord.Longitude, coord.Latitude, &redis.GeoRadiusQuery{
		Radius:    radiusKM,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
		Count:     10,
		Sort:      "ASC",
	})
	locations, err := cmd.Result()
	if err != nil {
		return 0, err
	}
	log.Println("list of location: ", locations)
	if len(locations) == 0 {
		return 0, nil
	}
	shipperID, err := strconv.Atoi(locations[0].Name)
	if err != nil {
		return 0, err
	}
	return shipperID, nil

}
