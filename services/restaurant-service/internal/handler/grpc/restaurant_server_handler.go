package grpc

import (
	"context"
	"restaurant-service/internal/service"
	"restaurant-service/proto/pb"
)

type RestaurantServiceServer struct {
	pb.UnimplementedRestaurantServiceServer
	restaurantService service.RestaurantService
	menuItemService   service.MenuItemService
}

func NewRestaurantServiceServer(restaurantService service.RestaurantService, itemService service.MenuItemService) *RestaurantServiceServer {
	return &RestaurantServiceServer{
		restaurantService: restaurantService,
		menuItemService:   itemService,
	}
}

func (r *RestaurantServiceServer) GetRestaurantInfoByID(ctx context.Context, request *pb.GetRestaurantInfoByIDRequest) (*pb.GetRestaurantInfoByIDResponse, error) {
	restaurant, err := r.restaurantService.GetRestaurantInfoByID(ctx, uint(request.RestaurantId))
	if err != nil {
		return nil, err
	}
	return &pb.GetRestaurantInfoByIDResponse{
		Id:          uint32(restaurant.ID),
		UserId:      int32(restaurant.ManagerID),
		Name:        restaurant.Name,
		Description: restaurant.Description,
		Address:     restaurant.Address,
		PhoneNumber: restaurant.PhoneNumber,
		Email:       restaurant.Email,
		Status:      restaurant.Status,
	}, nil
}

func (r *RestaurantServiceServer) GetMenuItem(ctx context.Context, request *pb.GetMenuItemRequest) (*pb.GetMenuItemResponse, error) {
	menuItem, err := r.menuItemService.GetMenuItemByItemID(ctx, uint(request.MenuItemId))
	if err != nil {
		return nil, err
	}
	return &pb.GetMenuItemResponse{
		Id:           uint32(menuItem.ID),
		RestaurantId: uint32(menuItem.RestaurantID),
		Name:         menuItem.Description,
		Description:  menuItem.Description,
		Price:        menuItem.Price,
		Available:    menuItem.Available,
	}, nil
}
