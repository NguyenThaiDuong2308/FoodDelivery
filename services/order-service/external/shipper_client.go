package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ShipperStatusClient interface {
	UpdateShipperStatus(ctx context.Context, status string, shipperID int) error
}

type shipperStatusClient struct {
	client *http.Client
}

func NewShipperStatusClient(client *http.Client) ShipperStatusClient {
	return &shipperStatusClient{client: client}
}

func (s *shipperStatusClient) UpdateShipperStatus(ctx context.Context, status string, shipperID int) error {
	url := fmt.Sprintf("http://shipper-service:8082/shipper/%d", shipperID)
	body := map[string]string{
		"status": status,
	}
	jsonBody, err := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
