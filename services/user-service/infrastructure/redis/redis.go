package redis

import (
	"log"
	"user-service/config"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("Failed to connect to redis")
		return nil, err
	}
	log.Println("Connected to redis")
	return rdb, nil
}
