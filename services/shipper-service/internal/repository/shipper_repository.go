package repository

import (
	"context"
	"shipper-service/internal/models"

	"gorm.io/gorm"
)

type ShipperRepository interface {
	CreateShipper(ctx context.Context, shipper *models.Shippers) error
	GetShipperInfoByID(ctx context.Context, id uint) (*models.Shippers, error)
	UpdateShipper(ctx context.Context, shipper *models.Shippers) error
	GetShipperList(ctx context.Context) (*[]models.Shippers, error)
	DeleteShipper(ctx context.Context, id uint) error
	GetShipperByUserID(ctx context.Context, userID uint) (*models.Shippers, error)
}

type shipperRepository struct {
	db *gorm.DB
}

func NewShipperRepository(db *gorm.DB) ShipperRepository {
	return &shipperRepository{db: db}
}

func (r *shipperRepository) CreateShipper(ctx context.Context, shipper *models.Shippers) error {
	result := r.db.WithContext(ctx).Create(shipper)
	if result.Error != nil {
		return result.Error
	}
	return result.Error
}

func (r *shipperRepository) GetShipperInfoByID(ctx context.Context, id uint) (*models.Shippers, error) {
	var shipper models.Shippers
	result := r.db.WithContext(ctx).First(&shipper, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &shipper, nil
}

func (r *shipperRepository) UpdateShipper(ctx context.Context, shipper *models.Shippers) error {
	result := r.db.WithContext(ctx).Save(shipper)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *shipperRepository) GetShipperList(ctx context.Context) (*[]models.Shippers, error) {
	var shippers []models.Shippers
	results := r.db.WithContext(ctx).Find(&shippers)
	if results.Error != nil {
		return nil, results.Error
	}
	return &shippers, nil
}

func (r *shipperRepository) DeleteShipper(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Shippers{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *shipperRepository) GetShipperByUserID(ctx context.Context, userID uint) (*models.Shippers, error) {
	var shipper models.Shippers
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&shipper)
	if result.Error != nil {
		return nil, result.Error
	}
	return &shipper, nil
}
