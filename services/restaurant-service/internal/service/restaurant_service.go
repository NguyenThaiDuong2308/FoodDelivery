package service

import (
	"context"
	"restaurant-service/internal/models"
	"restaurant-service/internal/repository"
)

type RestaurantService interface {
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
	GetRestaurantInfoByID(ctx context.Context, id uint) (*models.Restaurant, error)
	UpdateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
	GetRestaurantList(ctx context.Context) (*[]models.Restaurant, error)
	DeleteRestaurant(ctx context.Context, id uint) error
	GetRestaurantByUserID(ctx context.Context, userID int) (*models.Restaurant, error)
}

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
	}
}

func (s *restaurantService) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	return s.restaurantRepo.CreateRestaurant(ctx, restaurant)
}

func (s *restaurantService) GetRestaurantInfoByID(ctx context.Context, id uint) (*models.Restaurant, error) {
	return s.restaurantRepo.GetRestaurantInfoByID(ctx, id)
}

func (s *restaurantService) UpdateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	return s.restaurantRepo.UpdateRestaurant(ctx, restaurant)
}

func (s *restaurantService) GetRestaurantList(ctx context.Context) (*[]models.Restaurant, error) {
	return s.restaurantRepo.GetRestaurantList(ctx)
}

func (s *restaurantService) DeleteRestaurant(ctx context.Context, id uint) error {
	return s.restaurantRepo.DeleteRestaurant(ctx, id)
}

func (s *restaurantService) GetRestaurantByUserID(ctx context.Context, userID int) (*models.Restaurant, error) {
	return s.restaurantRepo.GetRestaurantByUserID(ctx, userID)
}
