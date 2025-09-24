package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RestaurantLocationClient interface {
	GetRestaurantLocation(ctx context.Context, restaurantID int) (string, error)
}
type restaurantLocationClient struct {
	client *http.Client
}

func NewRestaurantLocationClient(client *http.Client) RestaurantLocationClient {
	return &restaurantLocationClient{client: client}
}

type RestaurantLocationResponse struct {
	Location string `json:"location"`
}

func (r *restaurantLocationClient) GetRestaurantLocation(ctx context.Context, restaurantID int) (string, error) {
	url := fmt.Sprintf("http://restaurant-service:8081/restaurant/%d", restaurantID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	location, ok := data["address"].(string)
	if !ok {
		return "", errors.New("location field not found in response")
	}
	return location, nil
}
