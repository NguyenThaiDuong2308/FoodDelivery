package main

import (
	"log"
	"net/http"
	"shipper-assignment-service/api/route"
	"shipper-assignment-service/config"
	"shipper-assignment-service/external"
	"shipper-assignment-service/infrastructure/redis"
	"shipper-assignment-service/internal/handlers"
	"shipper-assignment-service/internal/repository"
	"shipper-assignment-service/internal/service"
	"time"

	"github.com/gin-gonic/gin"
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
	restaurantLocationClient := external.NewRestaurantLocationClient(client)
	shipperAssignmentService := service.NewShipperAssignmentService(geoRepo, locationRepo, restaurantLocationClient, mapboxClient, shipperStatusClient)
	shipperAssignmentHandler := handlers.NewShipperAssignmentHandler(shipperAssignmentService)
	r := gin.Default()
	route.SetupRoutes(r, shipperAssignmentHandler)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
