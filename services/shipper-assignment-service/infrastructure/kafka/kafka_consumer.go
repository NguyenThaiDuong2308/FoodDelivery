package kafka

import (
	"shipper-assignment-service/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	brokers []string
	group   string
	topic   string
	reader  *kafka.Reader
}

func NewKafka(cfg *config.Config) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: cfg.KafkaBrokers,
		group:   cfg.KafkaGroupID,
		topic:   cfg.KafkaTopic,
	}
}

func (k *KafkaConsumer) ConnectConsumer() *kafka.Reader {
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.brokers,
		GroupID:     k.group,
		Topic:       k.topic,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset,
		MaxWait:     1 * time.Second,
	})
	return k.reader
}
