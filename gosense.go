// Package gosense provides a generic sensor data generation engine.
// It supports various seeders, sensor functions, and publishers for generating realistic sensor data.
package gosense

import (
	"time"

	"github.com/Utsav-pixel/gosense/internal/engine"
)

// Re-export public types from internal/engine
type (
	// SensorData represents any sensor reading with generic data
	SensorData[T any] = engine.SensorData[T]

	// Quality represents the quality of sensor data
	Quality = engine.Quality

	// Seeder generates input values for sensor functions
	Seeder = engine.Seeder

	// SensorFunction defines the interface for sensor data generation functions
	SensorFunction[T any] = engine.SensorFunction[T]

	// Publisher defines the interface for publishing sensor data
	Publisher[T any] = engine.Publisher[T]

	// Config holds the engine configuration
	Config = engine.Config

	// Engine is the generic sensor engine
	Engine[T any] = engine.Engine[T]
)

// Re-export quality constants
const (
	QualityOK      = engine.QualityOK
	QualityNoisy   = engine.QualityNoisy
	QualityPartial = engine.QualityPartial
	QualityCorrupt = engine.QualityCorrupt
)

// Re-export seeder types and constructors
type (
	// TimeSeeder generates values based on time
	TimeSeeder = engine.TimeSeeder

	// RandomSeeder generates random values within a range
	RandomSeeder = engine.RandomSeeder

	// LinearSeeder generates values that increase linearly over time
	LinearSeeder = engine.LinearSeeder

	// CustomSeeder allows for custom generation functions
	CustomSeeder = engine.CustomSeeder

	// NormalSeeder generates values from a normal distribution
	NormalSeeder = engine.NormalSeeder
)

// NewTimeSeeder creates a new time-based seeder
func NewTimeSeeder(amplitude, frequency, offset float64) *TimeSeeder {
	return engine.NewTimeSeeder(amplitude, frequency, offset)
}

// NewRandomSeeder creates a new random seeder
func NewRandomSeeder(min, max float64) *RandomSeeder {
	return engine.NewRandomSeeder(min, max)
}

// NewLinearSeeder creates a new linear seeder
func NewLinearSeeder(slope, offset float64) *LinearSeeder {
	return engine.NewLinearSeeder(slope, offset)
}

// NewCustomSeeder creates a new custom seeder
func NewCustomSeeder(generateFunc func() float64) *CustomSeeder {
	return engine.NewCustomSeeder(generateFunc)
}

// NewNormalSeeder creates a new normal distribution seeder
func NewNormalSeeder(mean, stdDev float64) *NormalSeeder {
	return engine.NewNormalSeeder(mean, stdDev)
}

// Re-export sensor function types and constructors
type (
	// BasicSensorFunction provides a simple implementation for basic sensor data generation
	BasicSensorFunction[T any] = engine.BasicSensorFunction[T]

	// Function allows users to define their own sensor data generation logic
	Function[T any] = engine.Function[T]

	// LambdaSensorFunction provides a simple function wrapper for inline usage
	LambdaSensorFunction[T any] = engine.LambdaSensorFunction[T]
)

// NewBasicSensorFunction creates a new basic sensor function with a custom transform function
func NewBasicSensorFunction[T any](transformFunc func(float64, time.Time) T) *BasicSensorFunction[T] {
	return engine.NewBasicSensorFunction[T](transformFunc)
}

// NewFunction creates a new user-defined sensor function
func NewFunction[T any](generateFunc func(float64, time.Time) T) *Function[T] {
	return engine.NewFunction[T](generateFunc)
}

// NewLambdaSensorFunction creates a sensor function from a lambda/anonymous function
func NewLambdaSensorFunction[T any](lambda func(float64, time.Time) T) *LambdaSensorFunction[T] {
	return engine.NewLambdaSensorFunction[T](lambda)
}

// NewEngine creates a new generic sensor engine
func NewEngine[T any](
	config Config,
	seeder Seeder,
	function SensorFunction[T],
	publisher Publisher[T],
) *Engine[T] {
	return engine.NewEngine(config, seeder, function, publisher)
}

// Re-export configuration functions
func DefaultConfig() Config {
	return engine.DefaultConfig()
}

func HighThroughputConfig() Config {
	return engine.HighThroughputConfig()
}

func LowLatencyConfig() Config {
	return engine.LowLatencyConfig()
}
