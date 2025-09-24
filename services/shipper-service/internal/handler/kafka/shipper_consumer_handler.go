package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"shipper-service/internal/models"
	"shipper-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerHandler struct {
	service service.ShipperService
	reader  *kafka.Reader
}

func NewKafkaConsumerHandler(service service.ShipperService, reader *kafka.Reader) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		service: service,
		reader:  reader,
	}
}

func (h *KafkaConsumerHandler) StartConsume(ctx context.Context) error {
	for {
		msg, err := h.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return ctx.Err()
			}
			log.Printf("Error reading message: %s\n", err)
			continue
		}
		var receivedMessage struct {
			EventName string `json:"event_name"`
			UserID    uint   `json:"id"`
		}
		if err := json.Unmarshal(msg.Value, &receivedMessage); err != nil {
			log.Printf("Error unmarshalling message: %s\n", err)
			continue
		}
		if receivedMessage.EventName == "shipper_created" {
			shipper := models.Shippers{
				UserID: receivedMessage.UserID,
				Status: "offline",
			}
			if err := h.service.CreateShipper(ctx, &shipper); err != nil {
				log.Printf("Error creating shipper: %s\n", err)
			} else {
				log.Printf("Created shipper with ID: %d\n", shipper.ID)
			}
		} else if receivedMessage.EventName == "shipper_deleted" {
			shipper, err := h.service.GetShipperByUserID(ctx, receivedMessage.UserID)
			if err != nil {
				log.Printf("Error getting shipper: %s\n", err)
				continue
			}
			if err := h.service.DeleteShipper(ctx, shipper.ID); err != nil {
				log.Printf("Error deleting shipper: %s\n", err)
			} else {
				log.Printf("Deleted shipper with ID: %d\n", shipper.ID)
			}
		}
		if err := h.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %s\n", err)
		}
	}
}
