package repository

import (
	"context"
	"restaurant-service/internal/models"

	"gorm.io/gorm"
)

type MenuItemRepository interface {
	CreateMenuItem(ctx context.Context, menu *models.MenuItem) (*models.MenuItem, error)
	GetMenuListByRestaurantID(ctx context.Context, id uint) (*[]models.MenuItem, error)
	UpdateMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id uint) error
	GetAllAvailableMenuItem(ctx context.Context) (*[]models.MenuItem, error)
	GetAvailableMenuListByRestaurantID(ctx context.Context, id uint) (*[]models.MenuItem, error)
	GetMenuItemByItemID(ctx context.Context, id uint) (*models.MenuItem, error)
}
type menuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) CreateMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error) {
	result := r.db.WithContext(ctx).Create(&menuItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return menuItem, nil
}

func (r *menuItemRepository) GetMenuListByRestaurantID(ctx context.Context, id uint) (*[]models.MenuItem, error) {
	var menuItems []models.MenuItem
	result := r.db.WithContext(ctx).Where("restaurant_id = ?", id).Find(&menuItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menuItems, nil
}

func (r *menuItemRepository) GetAllAvailableMenuItem(ctx context.Context) (*[]models.MenuItem, error) {
	var menuItems []models.MenuItem
	result := r.db.WithContext(ctx).Where("available = true").Find(&menuItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menuItems, nil
}

func (r *menuItemRepository) GetAvailableMenuListByRestaurantID(ctx context.Context, id uint) (*[]models.MenuItem, error) {
	var menuItems []models.MenuItem
	result := r.db.WithContext(ctx).Where("restaurant_id = ? AND available = true", id).Find(&menuItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menuItems, nil
}

func (r *menuItemRepository) UpdateMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error) {
	result := r.db.WithContext(ctx).Save(&menuItem)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return menuItem, nil
}

func (r *menuItemRepository) DeleteMenuItem(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.MenuItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *menuItemRepository) GetMenuItemByItemID(ctx context.Context, id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	result := r.db.WithContext(ctx).First(&menuItem, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menuItem, nil
}
