package config

import (
	"os"
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
	KafkaTopic    string
	KafkaBrokers  []string
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFrom      string
}

var (
	cfg *Config
)

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
		KafkaBrokers:  []string{os.Getenv("KAFKA_BROKERS")},
		KafkaTopic:    os.Getenv("KAFKA_TOPIC"),
		SMTPHost:      os.Getenv("SMTP_Host"),
		SMTPPort:      os.Getenv("SMTP_Port"),
		SMTPUsername:  os.Getenv("SMTP_Username"),
		SMTPPassword:  os.Getenv("SMTP_Password"),
		SMTPFrom:      os.Getenv("SMTP_From"),
	}
	return cfg
}
