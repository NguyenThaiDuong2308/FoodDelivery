package main

import (
	"log"
	"user-service/api/routes"
	"user-service/config"
	"user-service/infrastructure/databases"
	"user-service/infrastructure/kafka"
	"user-service/infrastructure/redis"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.Println(cfg)
	db, err := databases.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kafkaProducer, err := kafka.NewKafka(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kafkaProducer.ConnectProducer()

	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//userRepo
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(redisClient)
	resetTokenRepo := repository.NewResetTokenRepository(db)
	userService := service.NewUserService(userRepo, kafkaProducer)
	mailer := service.NewSMTPMailer(cfg)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, resetTokenRepo, kafkaProducer, mailer)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	routes.SetupRoutes(r, authHandler, userHandler)

	r.Run(":" + cfg.ServerPort)
}
