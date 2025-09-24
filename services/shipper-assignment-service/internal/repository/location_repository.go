package repository

import (
	"context"
	"errors"
	"shipper-assignment-service/internal/models"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type LocationRepository interface {
	GetLocation(ctx context.Context, shipperID uint) (*models.Coordinate, error)
}
type locationRepository struct {
	rdb *redis.Client
}

const (
	AvailableLocationKey = "shipper_available:shipper_id:"
)

func NewLocationRepository(rdb *redis.Client) LocationRepository {
	return &locationRepository{rdb}
}

func (r *locationRepository) GetLocation(ctx context.Context, shipperID uint) (*models.Coordinate, error) {
	key := AvailableLocationKey
	pos, err := r.rdb.GeoPos(ctx, key, strconv.FormatUint(uint64(shipperID), 10)).Result()
	if err != nil {
		return nil, err
	}
	if len(pos) == 0 || pos[0] == nil {
		return nil, errors.New("location not found")
	}
	coord := &models.Coordinate{
		Longitude: pos[0].Longitude,
		Latitude:  pos[0].Latitude,
	}
	return coord, nil
}
