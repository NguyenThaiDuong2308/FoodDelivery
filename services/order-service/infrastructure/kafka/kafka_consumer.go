package kafka

import (
	"order-service/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	topic   string
	brokers []string
}

func NewKafkaConsumer(cfg *config.Config) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: cfg.KafkaBrokers,
		topic:   cfg.KafkaConsumeTopic,
	}
}
func (k *KafkaConsumer) ConnectConsumer() *kafka.Reader {
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.brokers,
		Topic:       k.topic,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset,
		MaxWait:     1 * time.Second,
	})
	return k.reader
}
