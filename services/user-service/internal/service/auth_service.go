package service

import (
	"context"
	"errors"
	"log"
	"time"
	"user-service/infrastructure/kafka"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, user models.User) (*models.User, error)
	Login(ctx context.Context, email string, password string) (string, string, error)
	ResetPassword(ctx context.Context, email string, oldPassword string, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, userID uint) error
	SendResetPasswordEmail(ctx context.Context, email string) error
	ResetForgotPassword(ctx context.Context, token string, newPassword string) error
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	resetTokenRepo   repository.ResetTokenRepository
	kafkaProducer    *kafka.KafkaProducer
	mailer           Mailer
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, resetTokenRepo repository.ResetTokenRepository, producer *kafka.KafkaProducer, mailer Mailer) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		resetTokenRepo:   resetTokenRepo,
		kafkaProducer:    producer,
		mailer:           mailer,
	}
}

const CreatedShipperEvent = "shipper_created"

//const Resetink = "abc"

func (s *authService) Register(ctx context.Context, user models.User) (*models.User, error) {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hash
	if user.Role == AdminRole {
		return nil, errors.New("you are not allowed to create admin")
	}
	if err := s.userRepo.CreateUser(ctx, &user); err != nil {
		return nil, err
	} else {
		if user.Role == ShipperRole {
			event := kafka.UserEvent{
				EventName: CreatedShipperEvent,
				ID:        user.ID,
			}
			log.Println("Created Shipper")
			if err := s.kafkaProducer.PublishUserEvent(ctx, event); err != nil {
				return nil, err
			}
		}

		return &user, nil
	}

}

func (s *authService) Login(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("user not found")
	}

	if utils.CheckPasswordHash(password, user.Password) == false {
		return "", "", errors.New("Invalid password")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Role, email, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Role, email, 24*time.Hour*7)
	if err != nil {
		return "", "", err
	}

	err = s.refreshTokenRepo.SetRefreshToken(ctx, user.ID, refreshToken, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}
	userIDRaw, ok := claims["user_id"]
	if !ok {
		return "", errors.New("not found user_id")
	}
	userID, ok := userIDRaw.(float64)
	if !ok {
		return "", errors.New("wrong type of userID")
	}
	storedToken, err := s.refreshTokenRepo.GetRefreshToken(ctx, uint(userID))
	if err != nil {
		return "", err
	}
	if storedToken != refreshToken {
		return "", errors.New("invalid token")
	}
	user, err := s.userRepo.GetUserById(ctx, uint(userID))
	if err != nil {
		return "", err
	}
	newAccessToken, err := utils.GenerateToken(user.ID, user.Role, user.Email, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	err := s.refreshTokenRepo.DeleteRefreshToken(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, email string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil || utils.CheckPasswordHash(oldPassword, user.Password) == false {
		return errors.New("invalid user or password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) SendResetPasswordEmail(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil
	}
	resetToken := models.PasswordResetToken{
		ID:         uuid.New(),
		Token:      uuid.New().String(),
		UserID:     user.ID,
		Expires_at: time.Now().Add(15 * time.Minute),
	}
	if err := s.resetTokenRepo.CreateResetToken(ctx, &resetToken); err != nil {
		return err
	}
	token := resetToken.Token
	err = s.mailer.SendResetPasswordEmail(ctx, email, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) ResetForgotPassword(ctx context.Context, token string, newPassword string) error {
	tokenData, err := s.resetTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	err = s.resetTokenRepo.MarkUsedToken(ctx, tokenData.ID)
	if err != nil {
		return err
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user, err := s.userRepo.GetUserById(ctx, tokenData.UserID)
	if err != nil {
		return err
	}
	user.Password = hash
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}
