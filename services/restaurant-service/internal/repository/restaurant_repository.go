package repository

import (
	"context"
	"restaurant-service/internal/models"

	"gorm.io/gorm"
)

type RestaurantRepository interface {
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
	GetRestaurantInfoByID(ctx context.Context, id uint) (*models.Restaurant, error)
	UpdateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
	GetRestaurantList(ctx context.Context) (*[]models.Restaurant, error)
	GetOpenRestaurantList(ctx context.Context) (*[]models.Restaurant, error)
	DeleteRestaurant(ctx context.Context, id uint) error
	GetRestaurantByManagerID(ctx context.Context, managerID int) (*models.Restaurant, error)
}

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	result := r.db.WithContext(ctx).Create(restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *restaurantRepository) GetRestaurantInfoByID(ctx context.Context, id uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	result := r.db.WithContext(ctx).Model(models.Restaurant{}).First(&restaurant, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}

func (r *restaurantRepository) UpdateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	result := r.db.WithContext(ctx).Model(models.Restaurant{}).Save(restaurant)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *restaurantRepository) GetRestaurantList(ctx context.Context) (*[]models.Restaurant, error) {
	var restaurants []models.Restaurant
	results := r.db.WithContext(ctx).Model(models.Restaurant{}).Find(&restaurants)
	if results.Error != nil {
		return nil, results.Error
	}
	return &restaurants, nil
}

func (r *restaurantRepository) GetOpenRestaurantList(ctx context.Context) (*[]models.Restaurant, error) {
	var restaurants []models.Restaurant
	results := r.db.WithContext(ctx).Model(models.Restaurant{}).Where("status = open").Find(&restaurants)
	if results.Error != nil {
		return nil, results.Error
	}
	return &restaurants, nil
}
func (r *restaurantRepository) DeleteRestaurant(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(models.Restaurant{}).Delete(&models.Restaurant{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *restaurantRepository) GetRestaurantByManagerID(ctx context.Context, managerID int) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	result := r.db.WithContext(ctx).Model(models.Restaurant{}).Where("manager_id = ?", managerID).First(&restaurant)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}
