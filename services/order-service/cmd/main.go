package main

import (
	"context"
	"log"
	"net/http"
	"order-service/api/route"
	"order-service/config"
	"order-service/external"
	"order-service/infrastructure/databases"
	"order-service/infrastructure/kafka"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/service"
	"time"

	"github.com/gin-gonic/gin"
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
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	shipperClient := external.NewShipperStatusClient(client)
	orderService := service.NewOrderService(orderRepo, restaurantClient, kafkaProducer, shipperClient)
	kafkaConsumer := kafka.NewKafkaConsumer(cfg)
	reader := kafkaConsumer.ConnectConsumer()
	orderHandler := handler.NewOrderHandler(orderService)
	consumerHandler := handler.NewKafkaConsumer(orderService, reader)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := consumerHandler.StartConsume(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	r := gin.Default()
	route.SetupRoutes(r, orderHandler)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}

}
