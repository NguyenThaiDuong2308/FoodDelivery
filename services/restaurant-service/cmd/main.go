package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"restaurant-service/api/route"
	"restaurant-service/config"
	"restaurant-service/infrastructure/databases"
	"restaurant-service/infrastructure/kafka"
	"restaurant-service/internal/handler"
	kafka2 "restaurant-service/internal/handler/kafka"
	"restaurant-service/internal/repository"
	"restaurant-service/internal/service"
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
	restaurantRepo := repository.NewRestaurantRepository(db)
	menuItemRepo := repository.NewMenuItemRepository(db)
	restaurantService := service.NewRestaurantService(restaurantRepo)
	menuItemService := service.NewMenuItemService(menuItemRepo)

	kafkaConsumer := kafka.NewKafka(cfg)
	reader := kafkaConsumer.ConnectConsumer()

	consumerHandler := kafka2.NewKafkaConsumerHandler(restaurantService, reader)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := consumerHandler.StartConsume(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	restaurantHandler := handler.NewRestaurantHandler(restaurantService)
	menuItemHandler := handler.NewMenuItemHandler(menuItemService, restaurantService)

	r := gin.Default()
	route.SetupRoutes(r, restaurantHandler, menuItemHandler)

	go func() {
		if err := r.Run(":" + cfg.ServerPort); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down service...")
	cancel()
	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}
}
