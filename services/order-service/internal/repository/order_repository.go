package repository

import (
	"context"
	"order-service/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	GetOrderByCustomerID(ctx context.Context, customerID int) (*[]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) error
	AssignShipperOrder(ctx context.Context, order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) AssignShipperOrder(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Save(order)
	return result.Error
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	result := r.db.WithContext(ctx).Preload("OrderItems").First(&order, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (r *orderRepository) GetOrderByCustomerID(ctx context.Context, customerID int) (*[]models.Order, error) {
	var orders []models.Order
	result := r.db.WithContext(ctx).Preload("OrderItems").Where("customer_id = ?", customerID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orders, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	var order models.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(&order).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
