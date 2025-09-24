package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"shipper-assignment-service/internal/models"
)

type MapboxResponse struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

type MapboxClient interface {
	ConvertCoordinate(address string) (*models.Coordinate, error)
	CalculateDistance(start *models.Coordinate, end *models.Coordinate) (float64, error)
}

type mapboxClient struct {
	mapboxToken string
	client      *http.Client
}

func NewMapboxClient(mapboxToken string, client *http.Client) MapboxClient {
	return &mapboxClient{
		mapboxToken: mapboxToken,
		client:      client,
	}
}

func (m *mapboxClient) ConvertCoordinate(address string) (*models.Coordinate, error) {
	log.Println(address)
	encodedAddress := url.QueryEscape(address)
	apiURL := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s", encodedAddress, m.mapboxToken)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mapbox responded with %s", resp.Status)
	}
	var response MapboxResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if len(response.Features) == 0 {
		return nil, errors.New("Not found any location")
	}
	coordinates := response.Features[0].Geometry.Coordinates
	if len(coordinates) < 2 {
		return nil, errors.New("Invalid coordinates")
	}
	log.Println("list coordinates:", coordinates)
	coordinate := &models.Coordinate{
		Longitude: coordinates[0],
		Latitude:  coordinates[1],
	}
	return coordinate, nil
}

type mapboxResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
	} `json:"routes"`
}

func (m *mapboxClient) CalculateDistance(start *models.Coordinate, end *models.Coordinate) (float64, error) {
	url := fmt.Sprintf("https://api.mapbox.com/directions/v5/mapbox/driving/%f,%f;%f,%f?access_token=%s",
		start.Longitude, start.Latitude, end.Longitude, end.Latitude, m.mapboxToken)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("mapbox responded with %s", resp.Status)
	}
	var response mapboxResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}
	if len(response.Routes) == 0 {
		return 0, errors.New("Not found any routes")
	}
	return response.Routes[0].Distance, nil
}
