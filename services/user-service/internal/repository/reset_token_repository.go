package repository

import (
	"context"
	"time"
	"user-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResetTokenRepository interface {
	CreateResetToken(ctx context.Context, token *models.PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	MarkUsedToken(ctx context.Context, id uuid.UUID) error
	DeleteExpiredToken(ctx context.Context) error
}

type resetTokenRepository struct {
	db *gorm.DB
}

func NewResetTokenRepository(db *gorm.DB) ResetTokenRepository {
	return &resetTokenRepository{
		db: db,
	}
}
func (r *resetTokenRepository) CreateResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *resetTokenRepository) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.WithContext(ctx).Where("token = ? AND used = FALSE AND expires_at > ?", token, time.Now()).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *resetTokenRepository) MarkUsedToken(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&models.PasswordResetToken{}).Where("id = ?", id).Update("used", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *resetTokenRepository) DeleteExpiredToken(ctx context.Context) error {
	err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
