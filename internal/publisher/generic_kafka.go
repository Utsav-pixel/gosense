package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
	"github.com/segmentio/kafka-go"
)

// GenericKafkaPublisher is a generic Kafka publisher
type GenericKafkaPublisher[T any] struct {
	writer *kafka.Writer
	batch  []kafka.Message
	mutex  sync.Mutex
}

// NewGenericKafkaPublisher creates a new generic Kafka publisher
func NewGenericKafkaPublisher[T any](brokers []string, topic string) *GenericKafkaPublisher[T] {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	})
	return &GenericKafkaPublisher[T]{
		writer: writer,
		batch:  make([]kafka.Message, 0, 100),
	}
}

// Publish publishes a single sensor data point
func (k *GenericKafkaPublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(data.ID),
		Value: value,
		Time:  time.Now(),
	}
	return k.writer.WriteMessages(ctx, msg)
}

// PublishBatch publishes a batch of sensor data points
func (k *GenericKafkaPublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	messages := make([]kafka.Message, len(data))
	for i, d := range data {
		value, err := json.Marshal(d)
		if err != nil {
			return err
		}
		messages[i] = kafka.Message{
			Key:   []byte(d.ID),
			Value: value,
			Time:  time.Now(),
		}
	}
	return k.writer.WriteMessages(ctx, messages...)
}

// Close closes the Kafka publisher
func (k *GenericKafkaPublisher[T]) Close() error {
	fmt.Println("Closing Kafka publisher")
	return k.writer.Close()
}
