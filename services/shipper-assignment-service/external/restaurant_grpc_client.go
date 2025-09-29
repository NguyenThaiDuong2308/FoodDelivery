package external

import (
	"context"
	"fmt"
	"log"
	"shipper-assignment-service/config"
	"shipper-assignment-service/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGrpcClient(cfg *config.Config) (client pb.RestaurantServiceClient, err error) {
	address := fmt.Sprintf("restaurant-service:%s", cfg.GRPCServerPort)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client = pb.NewRestaurantServiceClient(conn)
	return client, nil
}

type RestaurantClient interface {
	GetRestaurantLocationByID(ctx context.Context, id uint32) (string, error)
}

type restaurantClient struct {
	client pb.RestaurantServiceClient
}

func NewRestaurantClient(client pb.RestaurantServiceClient) RestaurantClient {
	return &restaurantClient{client: client}
}

func (c *restaurantClient) GetRestaurantLocationByID(ctx context.Context, id uint32) (string, error) {
	resp, err := c.client.GetRestaurantInfoByID(ctx, &pb.GetRestaurantInfoByIDRequest{
		RestaurantId: id,
	})
	if err != nil {
		return "", err
	}
	log.Println("resp: ", resp)
	return resp.Address, nil
}
