package service

import (
	"context"
	"shipper-service/internal/models"
	"shipper-service/internal/repository"
)

type ShipperService interface {
	CreateShipper(ctx context.Context, shipper *models.Shippers) error
	GetShipperInfoByID(ctx context.Context, id uint) (*models.Shippers, error)
	UpdateShipperStatus(ctx context.Context, shipper *models.Shippers) error
	UpdateShipperLocation(ctx context.Context, shipperID uint, long float64, lat float64) error
	GetShipperLocation(ctx context.Context, shipperID uint) (string, error)
	GetShipperByUserID(ctx context.Context, userID uint) (*models.Shippers, error)
	DeleteShipper(ctx context.Context, id uint) error
	GetShipperList(ctx context.Context) (*[]models.Shippers, error)
}
type shipperService struct {
	shipperRepo  repository.ShipperRepository
	locationRepo repository.LocationRepository
}

func NewShipperService(shipperRepo repository.ShipperRepository, locationRepo repository.LocationRepository) ShipperService {
	return &shipperService{
		shipperRepo:  shipperRepo,
		locationRepo: locationRepo,
	}
}

func (s *shipperService) CreateShipper(ctx context.Context, shipper *models.Shippers) error {
	err := s.shipperRepo.CreateShipper(ctx, shipper)
	if err != nil {
		return err
	}
	createdShipper, err := s.shipperRepo.GetShipperByUserID(ctx, shipper.UserID)
	if err != nil {
		return err
	}
	return s.locationRepo.SetStatus(ctx, createdShipper.ID, createdShipper.Status)
}

func (s *shipperService) GetShipperInfoByID(ctx context.Context, id uint) (*models.Shippers, error) {
	return s.shipperRepo.GetShipperInfoByID(ctx, id)
}
func (s *shipperService) UpdateShipperStatus(ctx context.Context, shipper *models.Shippers) error {
	status, err := s.locationRepo.GetStatus(ctx, shipper.ID)
	if err != nil {
		return err
	}
	if status == "available" || status == "busy" {
		err = s.locationRepo.RemoveLocation(ctx, shipper.ID, status)
		if err != nil {
			return err
		}
	}
	err = s.shipperRepo.UpdateShipper(ctx, shipper)
	if err != nil {
		return err
	}

	return s.locationRepo.SetStatus(ctx, shipper.ID, shipper.Status)
}

func (s *shipperService) GetShipperByUserID(ctx context.Context, userID uint) (*models.Shippers, error) {
	return s.shipperRepo.GetShipperByUserID(ctx, userID)
}

func (s *shipperService) DeleteShipper(ctx context.Context, id uint) error {
	return s.shipperRepo.DeleteShipper(ctx, id)
}

func (s *shipperService) GetShipperList(ctx context.Context) (*[]models.Shippers, error) {
	return s.shipperRepo.GetShipperList(ctx)
}

func (s *shipperService) UpdateShipperLocation(ctx context.Context, shipperID uint, long float64, lat float64) error {
	err := s.locationRepo.SetLocation(ctx, shipperID, long, lat)
	if err != nil {
		return err
	}
	return nil
}

func (s *shipperService) GetShipperLocation(ctx context.Context, shipperID uint) (string, error) {
	location, err := s.locationRepo.GetLocation(ctx, shipperID)
	if err != nil {
		return "", err
	}
	return location, nil
}
