package kafka

import (
	"context"
	"encoding/json"
	"log"
	"user-service/config"

	"github.com/segmentio/kafka-go"
)

type UserEvent struct {
	EventName string `json:"event_name"`
	ID        uint   `json:"id"`
}

type KafkaProducer struct {
	brokers []string
	writer  *kafka.Writer
	topic   string
}

func NewKafka(cfg *config.Config) (*KafkaProducer, error) {
	return &KafkaProducer{
		brokers: cfg.KafkaBrokers,
		topic:   cfg.KafkaTopic,
	}, nil
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

func (k *KafkaProducer) PublishUserEvent(ctx context.Context, event UserEvent) error {
	data, err := json.Marshal(event)
	log.Println(data)
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
