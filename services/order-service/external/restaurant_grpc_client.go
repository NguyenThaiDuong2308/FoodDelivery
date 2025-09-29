package external

import (
	"context"
	"fmt"
	"log"
	"order-service/config"
	"order-service/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGrpcClient(cfg *config.Config) (client pb.RestaurantServiceClient, err error) {
	address := fmt.Sprintf("restaurant-service:%s", cfg.GRPCServerPort)
	log.Println(address)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client = pb.NewRestaurantServiceClient(conn)
	return client, nil
}

type RestaurantClient interface {
	GetRestaurantInfoByID(ctx context.Context, id uint32) (*pb.GetRestaurantInfoByIDResponse, error)
	GetMenuItem(ctx context.Context, restaurantID, menuItemID uint32) (*pb.GetMenuItemResponse, error)
}
type restaurantClient struct {
	client pb.RestaurantServiceClient
}

func NewRestaurantClient(client pb.RestaurantServiceClient) RestaurantClient {
	return &restaurantClient{client: client}
}

func (c *restaurantClient) GetRestaurantInfoByID(ctx context.Context, id uint32) (*pb.GetRestaurantInfoByIDResponse, error) {
	log.Println("Starting find Restaurant Info")
	resp, err := c.client.GetRestaurantInfoByID(ctx, &pb.GetRestaurantInfoByIDRequest{
		RestaurantId: id,
	})
	if err != nil {
		log.Println("Can't find restaurant info because:", err)
		return nil, err
	}
	log.Println("Find Restaurant Info")
	return resp, nil
}

func (c *restaurantClient) GetMenuItem(ctx context.Context, restaurantID, menuItemID uint32) (*pb.GetMenuItemResponse, error) {
	resp, err := c.client.GetMenuItem(ctx, &pb.GetMenuItemRequest{
		RestaurantId: restaurantID,
		MenuItemId:   menuItemID,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
