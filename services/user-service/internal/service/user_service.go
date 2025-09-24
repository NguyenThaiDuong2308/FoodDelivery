package service

import (
	"context"
	"user-service/infrastructure/kafka"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetAllUserInfo(ctx context.Context) ([]*models.User, error)
}

type userService struct {
	userRepo      repository.UserRepository
	kafkaProducer *kafka.KafkaProducer
}

func NewUserService(userRepo repository.UserRepository, producer *kafka.KafkaProducer) UserService {
	return &userService{
		userRepo:      userRepo,
		kafkaProducer: producer}
}

const RestaurantAdminRole = "restaurant_admin"
const ShipperRole = "shipper"
const AdminRole = "admin"
const DeleteRestaurantEvent = "restaurant_deleted"
const DeleteShipperEvent = "shipper_deleted"

func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	if _, err := s.userRepo.GetUserById(ctx, id); err != nil {
		return nil, err
	}
	return s.userRepo.GetUserById(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	err = s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if user.Role == RestaurantAdminRole {
		event := kafka.UserEvent{
			EventName: DeleteRestaurantEvent,
			ID:        id,
		}
		if err := s.kafkaProducer.PublishUserEvent(ctx, event); err != nil {
			return err
		}
	} else if user.Role == ShipperRole {
		event := kafka.UserEvent{
			EventName: DeleteShipperEvent,
			ID:        id,
		}
		if err := s.kafkaProducer.PublishUserEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (s *userService) GetAllUserInfo(ctx context.Context) ([]*models.User, error) {
	users, err := s.userRepo.GetAllUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
