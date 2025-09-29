package config

import (
	"os"
)

type Config struct {
	ServerPort     string
	RedisAddress   string
	RedisPassword  string
	KafkaBrokers   []string
	KafkaTopic     string
	KafkaGroupID   string
	JWTSecret      string
	MapboxToken    string
	GRPCServerPort string
}

var cfg *Config

func LoadConfig() *Config {
	cfg = &Config{
		ServerPort:     os.Getenv("SERVER_PORT"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		RedisAddress:   os.Getenv("REDIS_ADDRESS"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		KafkaBrokers:   []string{os.Getenv("KAFKA_BROKERS")},
		KafkaTopic:     os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID:   os.Getenv("KAFKA_GROUP_ID"),
		MapboxToken:    os.Getenv("MAPBOX_TOKEN"),
		GRPCServerPort: os.Getenv("GRPC_SERVER_PORT"),
	}
	return cfg
}
