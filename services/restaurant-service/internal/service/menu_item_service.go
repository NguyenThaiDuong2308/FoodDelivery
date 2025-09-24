package service

import (
	"context"
	"restaurant-service/internal/models"
	"restaurant-service/internal/repository"
)

type MenuItemService interface {
	CreateMenuItem(ctx context.Context, menu *models.MenuItem) (*models.MenuItem, error)
	GetMenuListByRestaurantInfo(ctx context.Context, id uint) (*[]models.MenuItem, error)
	UpdateMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id uint) error
	GetMenuItemByItemID(ctx context.Context, menuItemId uint) (*models.MenuItem, error)
}

type menuItemService struct {
	menuItemRepo repository.MenuItemRepository
}

func (s *menuItemService) GetMenuItemByItemID(ctx context.Context, menuItemId uint) (*models.MenuItem, error) {
	return s.menuItemRepo.GetMenuItemByItemID(ctx, menuItemId)
}

func NewMenuItemService(menuItemRepo repository.MenuItemRepository) MenuItemService {
	return &menuItemService{
		menuItemRepo: menuItemRepo,
	}
}
func (s *menuItemService) CreateMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error) {
	return s.menuItemRepo.CreateMenuItem(ctx, menuItem)
}

func (s *menuItemService) GetMenuListByRestaurantInfo(ctx context.Context, id uint) (*[]models.MenuItem, error) {
	return s.menuItemRepo.GetMenuListByRestaurantID(ctx, id)
}

func (s *menuItemService) UpdateMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error) {
	return s.menuItemRepo.UpdateMenuItem(ctx, menuItem)
}

func (s *menuItemService) DeleteMenuItem(ctx context.Context, id uint) error {
	return s.menuItemRepo.DeleteMenuItem(ctx, id)
}
