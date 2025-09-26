package config

import (
	"os"
	"strings"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	KafkaBrokers   []string
	KafkaTopic     string
	GRPCServerPort string
}

var cfg *Config

func LoadConfig() *Config {
	cfg = &Config{
		ServerPort:     os.Getenv("SERVER_PORT"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		KafkaBrokers:   strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopic:     os.Getenv("KAFKA_TOPIC"),
		GRPCServerPort: os.Getenv("GRPC_SERVER_PORT"),
	}
	return cfg
}
