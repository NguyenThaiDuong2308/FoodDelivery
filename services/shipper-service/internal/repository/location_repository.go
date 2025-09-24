package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type LocationRepository interface {
	SetLocation(ctx context.Context, shipperID uint, longitude float64, latitude float64) error
	GetLocation(ctx context.Context, shipperID uint) (string, error)
	SetStatus(ctx context.Context, shipperID uint, status string) error
	GetStatus(ctx context.Context, shipperID uint) (string, error)
	RemoveLocation(ctx context.Context, shipperID uint, status string) error
}
type locationRepository struct {
	rdb *redis.Client
}

const (
	AvailableLocationKey = "shipper_available:shipper_id:"
	BusyLocationKey      = "shipper_busy:shipper_id:"
	ShipperStatusKey     = "shipper_status:shipper_id:"
)

func NewLocationRepository(rdb *redis.Client) LocationRepository {
	return &locationRepository{rdb}
}

func (r *locationRepository) SetLocation(ctx context.Context, shipperID uint, longitude float64, latitude float64) error {
	status, err := r.GetStatus(ctx, shipperID)
	if err != nil {
		return err
	}
	var key string
	if status == "available" {
		key = AvailableLocationKey
	} else if status == "busy" {
		key = BusyLocationKey
	} else if status == "offline" {
		_ = r.RemoveLocation(ctx, shipperID, "available")
		_ = r.RemoveLocation(ctx, shipperID, "busy")
		return nil
	} else {
		return errors.New("invalid status")
	}
	cmd := r.rdb.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      strconv.FormatUint(uint64(shipperID), 10),
		Longitude: longitude,
		Latitude:  latitude,
	})
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *locationRepository) GetLocation(ctx context.Context, shipperID uint) (string, error) {
	status, err := r.GetStatus(ctx, shipperID)
	if err != nil {
		return "", err
	}
	var key string
	if status == "available" {
		key = AvailableLocationKey
	} else if status == "busy" {
		key = BusyLocationKey
	} else if status == "offline" {
		return "", errors.New("shipper is offline")
	} else {
		return "", errors.New("invalid status")
	}
	pos, err := r.rdb.GeoPos(ctx, key, strconv.FormatUint(uint64(shipperID), 10)).Result()
	if err != nil {
		return "", err
	}
	if len(pos) == 0 || pos[0] == nil {
		return "", errors.New("location not found")
	}
	return fmt.Sprintf("POINT(%f %f)", pos[0].Longitude, pos[0].Latitude), nil
}

func (r *locationRepository) RemoveLocation(ctx context.Context, shipperID uint, status string) error {
	var key string
	if status == "available" {
		key = AvailableLocationKey
	} else if status == "busy" {
		key = BusyLocationKey
	} else {
		return errors.New("invalid status")
	}
	member := strconv.FormatUint(uint64(shipperID), 10)
	_, err := r.rdb.ZRem(ctx, key, member).Result()
	return err
}

func (r *locationRepository) SetStatus(ctx context.Context, shipperID uint, status string) error {
	key := ShipperStatusKey + strconv.FormatUint(uint64(shipperID), 10)
	return r.rdb.Set(ctx, key, status, 0).Err()
}

func (r *locationRepository) GetStatus(ctx context.Context, shipperID uint) (string, error) {
	key := ShipperStatusKey + strconv.FormatUint(uint64(shipperID), 10)
	status, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}
