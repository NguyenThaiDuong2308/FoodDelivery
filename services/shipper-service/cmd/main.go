package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"shipper-service/api/route"
	"shipper-service/config"
	"shipper-service/infrastructure/databases"
	"shipper-service/infrastructure/kafka"
	"shipper-service/infrastructure/redis"
	"shipper-service/internal/handler"
	kafka2 "shipper-service/internal/handler/kafka"
	"shipper-service/internal/repository"
	"shipper-service/internal/service"
	"syscall"

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
	shipperRepo := repository.NewShipperRepository(db)
	locationRepo := repository.NewLocationRepository(redisClient)
	shipperService := service.NewShipperService(shipperRepo, locationRepo)
	kafkaConsumer := kafka.NewKafka(cfg)
	reader := kafkaConsumer.ConnectConsumer()
	consumerHandler := kafka2.NewKafkaConsumerHandler(shipperService, reader)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := consumerHandler.StartConsume(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	shipperHandler := handler.NewShipperHandler(shipperService)
	r := gin.Default()
	route.SetupRoutes(r, shipperHandler)
	go func() {
		if err := r.Run(":" + cfg.ServerPort); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shipper service shutting down")
	cancel()
	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}

}
