package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"order-service/internal/dto"
	"order-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerHanlder struct {
	service service.OrderService
	reader  *kafka.Reader
}

func NewKafkaConsumer(service service.OrderService, reader *kafka.Reader) *KafkaConsumerHanlder {
	return &KafkaConsumerHanlder{
		service: service,
		reader:  reader,
	}
}

const AssignShipperEvent = "assign_shipper"

func (h *KafkaConsumerHanlder) StartConsume(ctx context.Context) error {
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
			EventName string  `json:"event_name"`
			OrderID   uint    `json:"order_id"`
			ShipperID int     `json:"shipper_id"`
			Distance  float64 `json:"distance"`
		}
		if err := json.Unmarshal(msg.Value, &receivedMessage); err != nil {
			log.Printf("Error unmarshalling message: %s\n", err)
			continue
		}
		if receivedMessage.EventName == AssignShipperEvent {
			event := dto.AssignShipperOrderRequest{
				OrderID:   receivedMessage.OrderID,
				ShipperID: receivedMessage.ShipperID,
				Distance:  receivedMessage.Distance,
			}
			log.Println("event consumer is:", event)
			err = h.service.AssignShipperOrder(ctx, event)
			if err != nil {
				log.Printf("Error assigning shipper order: %s\n", err)
				continue
			}
		}
		if err := h.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %s\n", err)
		}
	}
}
