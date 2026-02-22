package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
)

// ConsolePublisher publishes data to console for testing
type ConsolePublisher[T any] struct{}

func NewConsolePublisher[T any]() *ConsolePublisher[T] {
	return &ConsolePublisher[T]{}
}

func (c *ConsolePublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
	fmt.Printf("Single: ID=%s, Time=%v, Data=%+v, Quality=%s\n",
		data.ID, data.Timestamp.Format(time.RFC3339), data.Data, data.Quality)
	return nil
}

func (c *ConsolePublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
	fmt.Printf("Batch: %d items\n", len(data))
	for i, d := range data {
		fmt.Printf("  [%d] ID=%s, Time=%v, Data=%+v, Quality=%s\n",
			i, d.ID, d.Timestamp.Format(time.RFC3339), d.Data, d.Quality)
	}
	return nil
}

func (c *ConsolePublisher[T]) Close() error {
	fmt.Println("Console publisher closed")
	return nil
}

func main() {
	log.Println("Testing generic sensor engine with console output...")

	// Test with simple float64 data
	config := engine.DefaultConfig()
	config.ProductionRate = 500 * time.Millisecond
	config.BatchSize = 3
	config.BatchTimeout = 1 * time.Second

	// Simple time-based seeder
	seeder := engine.NewTimeSeeder(1.0, 0.1, 0.0)

	// Simple sensor function that generates temperature
	sensorFunc := engine.NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		// Input represents environmental factor (0-1)
		baseTemp := 20.0 + input*5.0
		// Add diurnal pattern
		hour := float64(timestamp.Hour()) + float64(timestamp.Minute())/60.0
		radian := (hour / 24.0) * 2 * math.Pi
		diurnal := 2.0 * math.Sin(radian-math.Pi/2)
		// Add noise
		noise := (rand.Float64() - 0.5) * 1.0
		return baseTemp + diurnal + noise
	})

	// Console publisher for testing
	publisher := NewConsolePublisher[float64]()

	// Create and start engine
	testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Starting test engine...")
	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}

	log.Println("Test completed successfully!")
}
