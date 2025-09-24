package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"restaurant-service/internal/service"

	"github.com/segmentio/kafka-go"
)

const RestaurantDeletedEvent = "restaurant_deleted"

type KafkaConsumerHandler struct {
	service service.RestaurantService
	reader  *kafka.Reader
}

func NewKafkaConsumerHandler(service service.RestaurantService, reader *kafka.Reader) *KafkaConsumerHandler {
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
			log.Println(err)
			continue
		}
		var receivedMessage struct {
			EventName string `json:"event_name"`
			UserID    uint   `json:"id"`
		}
		if err := json.Unmarshal(msg.Value, &receivedMessage); err != nil {
			log.Println(err)
			continue
		}
		if receivedMessage.EventName == RestaurantDeletedEvent {
			restaurant, err := h.service.GetRestaurantByUserID(ctx, int(receivedMessage.UserID))
			if err != nil {
				log.Println(err)
			}
			err = h.service.DeleteRestaurant(ctx, restaurant.ID)
			if err != nil {
				log.Println(err)
			}
		}
		if err := h.reader.CommitMessages(ctx, msg); err != nil {
			log.Println(err)
		}
	}
}
