package broker

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Host  string
	Port  string
	Topic string
}

func NewKafkaWriter(cfg *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Host + ":" + cfg.Port),
		Topic:                  cfg.Topic,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}
}

func ProduceNewUser(ctx context.Context, payload any, w *kafka.Writer) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal kafka message value")
	}
	err = w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: data,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to write kafka messages")
	}
}
