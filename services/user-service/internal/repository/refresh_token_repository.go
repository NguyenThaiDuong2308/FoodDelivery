package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepository interface {
	SetRefreshToken(ctx context.Context, userID uint, token string, expiresIn time.Duration) error
	GetRefreshToken(ctx context.Context, userID uint) (string, error)
	DeleteRefreshToken(ctx context.Context, userID uint) error
}

type tokenRepository struct {
	rdb *redis.Client
}

func NewRefreshTokenRepository(rdb *redis.Client) RefreshTokenRepository {
	return &tokenRepository{
		rdb: rdb,
	}
}

const RefreshTokenKey = "user:refresh_token:user_id"

func (t *tokenRepository) SetRefreshToken(ctx context.Context, userID uint, token string, expiresIn time.Duration) error {
	key := RefreshTokenKey + strconv.FormatUint(uint64(userID), 10)
	return t.rdb.Set(ctx, key, token, expiresIn).Err()
}

func (t *tokenRepository) GetRefreshToken(ctx context.Context, userID uint) (string, error) {
	key := RefreshTokenKey + strconv.FormatUint(uint64(userID), 10)
	token, err := t.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (t *tokenRepository) DeleteRefreshToken(ctx context.Context, userID uint) error {
	key := RefreshTokenKey + strconv.FormatUint(uint64(userID), 10)
	return t.rdb.Del(ctx, key).Err()
}
