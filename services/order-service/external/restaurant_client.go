package external

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"net/http"
//)
//
//type RestaurantClient interface {
//	GetMenuItem(ctx context.Context, restaurantID int, menuItemID int) (*MenuItemResponse, error)
//	GetRestaurantInfoByID(ctx context.Context, restaurantID int) (*RestaurantResponse, error)
//}
//
//type restaurantClient struct {
//	client *http.Client
//}
//
//func NewRestaurantClient(client *http.Client) RestaurantClient {
//	return &restaurantClient{client: client}
//}
//
//type MenuItemResponse struct {
//	ID           uint    `json:"id"`
//	RestaurantID uint    `json:"restaurant_id"`
//	Name         string  `json:"name"`
//	Description  string  `json:"description"`
//	Price        float64 `json:"price"`
//	Available    bool    `json:"available"`
//}
//
//func (r *restaurantClient) GetMenuItem(ctx context.Context, restaurantID int, menuItemID int) (*MenuItemResponse, error) {
//	url := fmt.Sprintf("http://restaurant-service:8081/restaurant/%d/menu/%d", restaurantID, menuItemID)
//	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := r.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//	var menuItem MenuItemResponse
//	err = json.NewDecoder(resp.Body).Decode(&menuItem)
//	if err != nil {
//		return nil, err
//	}
//	return &menuItem, nil
//}
//
//type RestaurantResponse struct {
//	Id     uint   `json:"id"`
//	UserId int    `json:"user_id"`
//	Status string `json:"status"`
//}
//
//func (r *restaurantClient) GetRestaurantInfoByID(ctx context.Context, restaurantID int) (*RestaurantResponse, error) {
//	url := fmt.Sprintf("http://restaurant-service:8081/restaurant/%d", restaurantID)
//	fmt.Println(url)
//	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := r.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//	var restaurant RestaurantResponse
//	err = json.NewDecoder(resp.Body).Decode(&restaurant)
//	if err != nil {
//		return nil, err
//	}
//	return &restaurant, nil
//}
