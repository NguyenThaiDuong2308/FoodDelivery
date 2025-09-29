package main

import (
	"context"
	"log"
	"net/http"
	"shipper-assignment-service/config"
	"shipper-assignment-service/external"
	"shipper-assignment-service/infrastructure/kafka"
	"shipper-assignment-service/infrastructure/redis"
	"shipper-assignment-service/internal/handlers"
	"shipper-assignment-service/internal/repository"
	"shipper-assignment-service/internal/service"
	"time"
)

const GeoKey = "shipper_available:shipper_id:"

func main() {
	cfg := config.LoadConfig()
	log.Println(cfg)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	locationRepo := repository.NewLocationRepository(redisClient)
	geoRepo := repository.NewGeoOperations(redisClient, GeoKey)
	mapboxClient := external.NewMapboxClient(cfg.MapboxToken, client)
	shipperStatusClient := external.NewShipperStatusClient(client)
	grpcConn, err := external.StartGrpcClient(cfg)
	restaurantLocationClient := external.NewRestaurantClient(grpcConn)
	if err != nil {
		log.Fatal(err)
	}
	//kafkaProducer := kafka.NewKafkaProducer(cfg)
	//kafkaProducer.ConnectProducer()
	//defer func() {
	//	if err := kafkaProducer.Close(); err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	shipperAssignmentService := service.NewShipperAssignmentService(geoRepo, locationRepo, restaurantLocationClient, mapboxClient, shipperStatusClient)
	kafkaConsumer := kafka.NewKafka(cfg)
	reader := kafkaConsumer.ConnectConsumer()
	consumerHandler := handlers.NewKafkaConsumerHandler(shipperAssignmentService, reader)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = consumerHandler.StartConsume(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//r := gin.Default()
	//route.SetupRoutes(r, shipperAssignmentHandler)
	//if err := r.Run(":" + cfg.ServerPort); err != nil {
	//	log.Fatal(err)
	//}
}
