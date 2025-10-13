package service

import (
	"context"
	"errors"
	"log"
	"order-service/external"
	"order-service/infrastructure/kafka"
	"order-service/internal/dto"
	"order-service/internal/models"
	"order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int, restaurantID int, orderItems []models.OrderItem) error
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	GetOrderByCustomerID(ctx context.Context, customerID int) (*[]models.Order, error)
	GetOrderByRestaurantID(ctx context.Context, restaurantID int) (*[]models.Order, error)
	GetOrderByShipperID(ctx context.Context, shipperID int) (*[]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) error
	AssignShipperOrder(ctx context.Context, request dto.AssignShipperOrderRequest) error
}

type orderService struct {
	orderRepo        repository.OrderRepository
	restaurantClient external.RestaurantClient
	kafkaProducer    *kafka.KafkaProducer
}

const OrderCreatedEvent = "order_created"
const CreatedOrderStatus = "created"
const CanceledOrderStatus = "cancelled"
const DeliveringOrderStatus = "delivering"

func (o *orderService) CreateOrder(ctx context.Context, customerID int, restaurantID int, orderItems []models.OrderItem) error {
	var order models.Order
	order.CustomerID = customerID
	order.RestaurantID = restaurantID
	order.OrderItems = orderItems
	order.ItemsPrice = 0
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
		order.ItemsPrice += menuItem.Price * float64(item.Quantity)
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

func (o *orderService) AssignShipperOrder(ctx context.Context, request dto.AssignShipperOrderRequest) error {
	order, err := o.orderRepo.GetOrderByID(ctx, request.OrderID)
	if err != nil {
		return err
	}
	log.Println("assigning shipper", request.ShipperID)
	if request.ShipperID == 0 {
		order.ShipperID = 0
		order.Status = CanceledOrderStatus
		order.DeliveryPrice = 0
		order.TotalPrice = order.DeliveryPrice + order.ItemsPrice
		log.Println(order)
		return o.orderRepo.AssignShipperOrder(ctx, order)
	}
	order.ShipperID = request.ShipperID
	order.Status = DeliveringOrderStatus
	order.DeliveryPrice = request.Distance * 0.05
	order.TotalPrice = order.DeliveryPrice + order.ItemsPrice
	return o.orderRepo.AssignShipperOrder(ctx, order)
}

func (o *orderService) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {

	return o.orderRepo.GetOrderByID(ctx, id)
}

func (o *orderService) GetOrderByCustomerID(ctx context.Context, customerID int) (*[]models.Order, error) {
	return o.orderRepo.GetOrderByCustomerID(ctx, customerID)
}

func (o *orderService) GetOrderByRestaurantID(ctx context.Context, restaurantID int) (*[]models.Order, error) {
	return o.orderRepo.GetOrderByRestaurantID(ctx, restaurantID)
}

func (o *orderService) GetOrderByShipperID(ctx context.Context, shipperID int) (*[]models.Order, error) {
	return o.orderRepo.GetOrderByShipperID(ctx, shipperID)
}

func (o *orderService) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	return o.orderRepo.UpdateOrderStatus(ctx, id, status)
}

func NewOrderService(orderRepo repository.OrderRepository, restaurantClient external.RestaurantClient, producer *kafka.KafkaProducer) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		restaurantClient: restaurantClient,
		kafkaProducer:    producer,
	}
}
