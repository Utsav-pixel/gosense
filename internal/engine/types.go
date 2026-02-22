package engine

import (
	"context"
	"time"
)

// SensorData represents any sensor reading with generic data
type SensorData[T any] struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
	Quality   Quality   `json:"quality"`
}

// Quality represents the quality of sensor data
type Quality string

const (
	QualityOK      Quality = "OK"
	QualityNoisy   Quality = "NOISY"
	QualityPartial Quality = "PARTIAL"
	QualityCorrupt Quality = "CORRUPT"
)

// Seeder generates input values for sensor functions
type Seeder interface {
	Generate() float64
}

// SensorFunction defines the interface for sensor data generation functions
type SensorFunction[T any] interface {
	Generate(input float64, timestamp time.Time) T
}

// Publisher defines the interface for publishing sensor data
type Publisher[T any] interface {
	Publish(ctx context.Context, data SensorData[T]) error
	PublishBatch(ctx context.Context, data []SensorData[T]) error
	Close() error
}

// Config holds the engine configuration
type Config struct {
	ProductionRate time.Duration // How often to generate data
	BatchSize      int           // Number of messages to batch together
	BatchTimeout   time.Duration // How long to wait before publishing a batch
	MaxWorkers     int           // Number of concurrent workers
}

// Engine is the generic sensor engine
type Engine[T any] struct {
	config    Config
	seeder    Seeder
	function  SensorFunction[T]
	publisher Publisher[T]
}

// NewEngine creates a new generic sensor engine
func NewEngine[T any](
	config Config,
	seeder Seeder,
	function SensorFunction[T],
	publisher Publisher[T],
) *Engine[T] {
	return &Engine[T]{
		config:    config,
		seeder:    seeder,
		function:  function,
		publisher: publisher,
	}
}
