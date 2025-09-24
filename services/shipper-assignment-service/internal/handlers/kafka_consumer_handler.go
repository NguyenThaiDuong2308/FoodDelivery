package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"shipper-assignment-service/internal/models"
	"shipper-assignment-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerHandler struct {
	service service.ShipperAssignmentService
	reader  *kafka.Reader
}

func NewKafkaConsumerHandler(service service.ShipperAssignmentService, reader *kafka.Reader) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		service: service,
		reader:  reader,
	}
}

func (h *KafkaConsumerHandler) StartConsume(ctx context.Context) (*models.Assignment, error) {
	for {
		msg, err := h.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("Error fetching message: %s\n", err)
				return nil, ctx.Err()
			}
			log.Println(err)
			continue
		}
		var receivedMessage struct {
			EventName    string `json:"event_name"`
			OrderID      uint   `json:"order_id"`
			RestaurantID int    `json:"restaurant_id"`
		}
		if err := json.Unmarshal(msg.Value, &receivedMessage); err != nil {
			log.Println(err)
			continue
		}
		if receivedMessage.EventName == "order_created" {
			event := models.OrderEvent{
				EventName:    receivedMessage.EventName,
				OrderID:      receivedMessage.OrderID,
				RestaurantID: receivedMessage.RestaurantID,
			}
			result, err := h.service.AssignNearestShipper(ctx, &event)
			if err != nil {
				log.Println(err)
				continue
			}
			return result, nil
		}
		if err := h.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %s\n", err)
		}
	}
}
