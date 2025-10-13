package service

import (
	"context"
	"log"
	"shipper-assignment-service/external"
	"shipper-assignment-service/internal/models"
	"shipper-assignment-service/internal/repository"
)

const BusyStatus = "busy"

type ShipperAssignmentService interface {
	AssignNearestShipper(ctx context.Context, event *models.OrderEvent) (*models.Assignment, error)
}
type shipperAssignmentService struct {
	geoRepo                  repository.GeoOperations
	locationRepo             repository.LocationRepository
	restaurantLocationClient external.RestaurantClient
	mapboxClient             external.MapboxClient
	shipperStatusClient      external.ShipperStatusClient
}

func NewShipperAssignmentService(geoRepo repository.GeoOperations, locationRepo repository.LocationRepository, restaurantLocationClient external.RestaurantClient, mapboxClient external.MapboxClient, shipperStatusClient external.ShipperStatusClient) ShipperAssignmentService {
	return &shipperAssignmentService{
		geoRepo:                  geoRepo,
		locationRepo:             locationRepo,
		restaurantLocationClient: restaurantLocationClient,
		mapboxClient:             mapboxClient,
		shipperStatusClient:      shipperStatusClient,
	}
}

func (s *shipperAssignmentService) AssignNearestShipper(ctx context.Context, event *models.OrderEvent) (*models.Assignment, error) {
	restaurantLocation, err := s.restaurantLocationClient.GetRestaurantLocationByID(ctx, uint32(event.RestaurantID))
	log.Println("location: ", restaurantLocation)
	if err != nil {
		return nil, err
	}
	restaurantCoord, err := s.mapboxClient.ConvertCoordinate(restaurantLocation)
	if err != nil {
		return nil, err
	}
	var result models.Assignment
	shipperID, err := s.geoRepo.FindNearestShipper(ctx, restaurantCoord, 100)
	if err != nil {
		return nil, err
	}
	result.ShipperID = shipperID
	result.OrderID = event.OrderID
	log.Println("this is shipperID ", shipperID)
	if shipperID == 0 {
		return &result, nil
	}
	shipperCoord, err := s.locationRepo.GetLocation(ctx, uint(shipperID))
	if err != nil {
		return nil, err
	}
	result.Distance, err = s.mapboxClient.CalculateDistance(restaurantCoord, shipperCoord)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("successful to calculate distance")
	err = s.shipperStatusClient.UpdateShipperStatus(ctx, BusyStatus, result.ShipperID)
	if err != nil {
		return nil, err
	}
	log.Println("successful to update status")
	return &result, nil
}
