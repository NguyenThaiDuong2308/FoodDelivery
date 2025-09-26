package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"order-service/api/route"
	"order-service/config"
	"order-service/external"
	"order-service/infrastructure/databases"
	"order-service/infrastructure/kafka"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	log.Println(cfg)
	db, err := databases.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}
	orderRepo := repository.NewOrderRepository(db)
	kafkaProducer := kafka.NewKafka(cfg)
	kafkaProducer.ConnectProducer()
	defer func() {
		if err := kafkaProducer.CloseProducer(); err != nil {
			log.Fatal(err)
		}
	}()
	grpcConn, err := external.StartGrpcClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	restaurantClient := external.NewRestaurantClient(grpcConn)
	orderService := service.NewOrderService(orderRepo, restaurantClient, kafkaProducer)
	orderHandler := handler.NewOrderHandler(orderService)
	r := gin.Default()
	route.SetupRoutes(r, orderHandler)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}

}
