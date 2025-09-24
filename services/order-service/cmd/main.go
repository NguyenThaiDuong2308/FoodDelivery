package main

import (
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
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	restaurantClient := external.NewRestaurantClient(client)
	kafkaProducer := kafka.NewKafka(cfg)
	kafkaProducer.ConnectProducer()
	defer func() {
		if err := kafkaProducer.CloseProducer(); err != nil {
			log.Fatal(err)
		}
	}()
	orderService := service.NewOrderService(orderRepo, restaurantClient, kafkaProducer)
	orderHandler := handler.NewOrderHandler(orderService)
	r := gin.Default()
	route.SetupRoutes(r, orderHandler)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}

}
