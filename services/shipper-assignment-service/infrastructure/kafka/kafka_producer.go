package kafka

import (
	"context"
	"encoding/json"
	"log"
	"shipper-assignment-service/config"

	"github.com/segmentio/kafka-go"
)

type AssignEvent struct {
	EventName string  `json:"event_name"`
	OrderID   uint    `json:"order_id"`
	ShipperID int     `json:"shipper_id"`
	Distance  float64 `json:"distance"`
}

type KafkaProducer struct {
	brokers []string
	writer  *kafka.Writer
	topic   string
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
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

func (k *KafkaProducer) PublishAssignEvent(ctx context.Context, event AssignEvent) error {
	data, err := json.Marshal(event)
	log.Println(string(data))
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

func (k *KafkaProducer) Close() error {
	if k.writer != nil {
		return k.writer.Close()
	}
	return nil
}
