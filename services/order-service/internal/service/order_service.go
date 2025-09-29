package service

import (
	"context"
	"errors"
	"log"
	"order-service/external"
	"order-service/infrastructure/kafka"
	"order-service/internal/models"
	"order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int, restaurantID int, orderItems []models.OrderItem) error
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	GetOrderByCustomerID(ctx context.Context, customerID int) (*[]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) error
}

type orderService struct {
	orderRepo        repository.OrderRepository
	restaurantClient external.RestaurantClient
	kafkaProducer    *kafka.KafkaProducer
}

const OrderCreatedEvent = "order_created"
const CreatedOrderStatus = "created"

func (o orderService) CreateOrder(ctx context.Context, customerID int, restaurantID int, orderItems []models.OrderItem) error {
	var order models.Order
	order.CustomerID = customerID
	order.RestaurantID = restaurantID
	order.OrderItems = orderItems
	order.TotalPrice = 0
	_, err := o.restaurantClient.GetRestaurantInfoByID(ctx, uint32(order.RestaurantID))
	if err != nil {
		return errors.New("Restaurant not found")
	}
	for _, item := range orderItems {
		menuItem, err := o.restaurantClient.GetMenuItem(ctx, uint32(order.RestaurantID), uint32(item.MenuItemID))
		if err != nil {
			return err
		}
		if menuItem == nil {
			return errors.New("MenuItem not found")
		}
		if menuItem.Available == false {
			return errors.New("MenuItem not available")
		}
		order.TotalPrice += menuItem.Price * float64(item.Quantity)
	}
	order.Status = CreatedOrderStatus
	if err := o.orderRepo.CreateOrder(ctx, &order); err != nil {
		return err
	} else {
		event := kafka.OrderEvent{
			EventName:    OrderCreatedEvent,
			OrderID:      order.ID,
			RestaurantID: order.RestaurantID,
		}
		log.Println("finding nearest shipper")
		if err := o.kafkaProducer.PublishUserEvent(ctx, event); err != nil {
			return err
		}
		return nil
	}
}

func (o orderService) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {

	return o.orderRepo.GetOrderByID(ctx, id)
}

func (o orderService) GetOrderByCustomerID(ctx context.Context, customerID int) (*[]models.Order, error) {
	return o.orderRepo.GetOrderByCustomerID(ctx, customerID)
}

func (o orderService) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	return o.orderRepo.UpdateOrderStatus(ctx, id, status)
}

func NewOrderService(orderRepo repository.OrderRepository, restaurantClient external.RestaurantClient, producer *kafka.KafkaProducer) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		restaurantClient: restaurantClient,
		kafkaProducer:    producer,
	}
}
