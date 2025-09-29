package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"restaurant-service/api/route"
	"restaurant-service/config"
	"restaurant-service/infrastructure/databases"
	"restaurant-service/infrastructure/kafka"
	grpc2 "restaurant-service/internal/handler/grpc"
	kafka2 "restaurant-service/internal/handler/kafka"
	"restaurant-service/internal/handler/rest"
	"restaurant-service/internal/repository"
	"restaurant-service/internal/service"
	"restaurant-service/proto/pb"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	restaurantHandler := rest.NewRestaurantHandler(restaurantService)
	menuItemHandler := rest.NewMenuItemHandler(menuItemService, restaurantService)

	r := gin.Default()
	route.SetupRoutes(r, restaurantHandler, menuItemHandler)

	go func() {
		if err := r.Run(":" + cfg.ServerPort); err != nil {
			log.Fatal(err)
		}
	}()
	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("Faled to listen gRPC port: %v", err)
	}
	grpcServer := grpc.NewServer()
	grpcHandler := grpc2.NewRestaurantServiceServer(restaurantService, menuItemService)

	pb.RegisterRestaurantServiceServer(grpcServer, grpcHandler)
	go func() {
		log.Printf("gRPC server started on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down service...")
	grpcServer.GracefulStop()
	cancel()
	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}
}
