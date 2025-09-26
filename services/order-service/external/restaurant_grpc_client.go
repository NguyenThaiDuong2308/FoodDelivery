package external

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"order-service/config"
	"proto/pb"
)

func StartGrpcClient(cfg *config.Config) (client pb.RestaurantServiceClient, err error) {
	conn, err := grpc.NewClient(cfg.GRPCServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	resp, err := c.client.GetRestaurantInfoByID(ctx, &pb.GetRestaurantInfoByIDRequest{
		RestaurantId: id,
	})
	if err != nil {
		return nil, err
	}
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
