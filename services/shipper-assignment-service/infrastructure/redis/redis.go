package redis

import (
	"context"
	"log"
	"shipper-assignment-service/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to redis ", err)
		return nil, err
	}
	log.Println("Connected to redis")
	return rdb, nil
}
