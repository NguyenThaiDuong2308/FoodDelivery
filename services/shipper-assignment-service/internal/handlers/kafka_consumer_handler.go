package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	kafka2 "shipper-assignment-service/infrastructure/kafka"
	"shipper-assignment-service/internal/models"
	"shipper-assignment-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerHandler struct {
	service service.ShipperAssignmentService
	reader  *kafka.Reader
	writer  *kafka2.KafkaProducer
}

func NewKafkaConsumerHandler(service service.ShipperAssignmentService, reader *kafka.Reader, writer *kafka2.KafkaProducer) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		service: service,
		reader:  reader,
		writer:  writer,
	}
}

const OrderCreatedEvent = "order_created"
const AssignShipperEvent = "assign_shipper"

func (h *KafkaConsumerHandler) StartConsume(ctx context.Context) error {
	for {
		msg, err := h.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("Error fetching message: %s\n", err)
				return ctx.Err()
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
		if receivedMessage.EventName == OrderCreatedEvent {
			event := models.OrderEvent{
				EventName:    receivedMessage.EventName,
				OrderID:      receivedMessage.OrderID,
				RestaurantID: receivedMessage.RestaurantID,
			}
			log.Println(event)
			result, err := h.service.AssignNearestShipper(ctx, &event)

			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("successful to handle message", result)
			assignEvent := kafka2.AssignEvent{
				EventName: AssignShipperEvent,
				OrderID:   result.OrderID,
				ShipperID: result.ShipperID,
				Distance:  result.Distance,
			}
			log.Println("assign event:", assignEvent)
			produceErr := h.writer.PublishAssignEvent(ctx, assignEvent)
			if produceErr != nil {
				log.Println("Error producing shipper_assigned event:", produceErr)
			}
		}
		if err := h.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %s\n", err)
		}
	}
}
