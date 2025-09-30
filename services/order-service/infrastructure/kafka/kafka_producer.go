package kafka

import (
	"context"
	"encoding/json"
	"order-service/config"

	"github.com/segmentio/kafka-go"
)

type OrderEvent struct {
	EventName    string `json:"event_name"`
	OrderID      uint   `json:"order_id"`
	RestaurantID int    `json:"restaurant_id"`
}

type KafkaProducer struct {
	brokers []string
	writer  *kafka.Writer
	topic   string
}

func NewKafka(cfg *config.Config) *KafkaProducer {
	return &KafkaProducer{
		brokers: cfg.KafkaBrokers,
		topic:   cfg.KafkaProduceTopic,
	}
}

func (k *KafkaProducer) ConnectProducer() *kafka.Writer {
	k.writer = &kafka.Writer{
		Addr:         kafka.TCP(k.brokers...),
		Topic:        k.topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	return k.writer
}

func (k *KafkaProducer) PublishUserEvent(ctx context.Context, event OrderEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = k.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *KafkaProducer) CloseProducer() error {
	if k.writer != nil {
		return k.writer.Close()
	}
	return nil
}
