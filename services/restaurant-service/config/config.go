package config

import (
	"os"
	"strings"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	RedisAddress  string
	RedisPassword string
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaGroupID  string
}

var cfg *Config

func LoadConfig() *Config {
	cfg = &Config{
		ServerPort:    os.Getenv("SERVER_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisAddress:  os.Getenv("REDIS_ADDRESS"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		KafkaBrokers:  strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopic:    os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID:  os.Getenv("KAFKA_GROUP_ID"),
	}
	return cfg

}
