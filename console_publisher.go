package gosense

import (
	"context"
	"fmt"
)

// ConsolePublisher is a simple publisher that outputs sensor data to the console
// Useful for development, testing, and demonstration purposes
type ConsolePublisher[T any] struct{}

// NewConsolePublisher creates a new console publisher
func NewConsolePublisher[T any]() *ConsolePublisher[T] {
	return &ConsolePublisher[T]{}
}

// Publish publishes a single sensor data point to the console
func (p *ConsolePublisher[T]) Publish(ctx context.Context, data SensorData[T]) error {
	fmt.Printf("📊 [%s] %+v\n", data.Quality, data.Data)
	return nil
}

// PublishBatch publishes a batch of sensor data points to the console
func (p *ConsolePublisher[T]) PublishBatch(ctx context.Context, data []SensorData[T]) error {
	fmt.Printf("📦 Batch of %d items:\n", len(data))
	for i, item := range data {
		fmt.Printf("  [%d] [%s] %+v\n", i, item.Quality, item.Data)
	}
	return nil
}

// Close closes the console publisher
func (p *ConsolePublisher[T]) Close() error {
	fmt.Println("🔚 Console publisher closed")
	return nil
}
