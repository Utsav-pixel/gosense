package engine

import (
	"context"
	"testing"
	"time"
)

func TestEngine_Integration_Float64(t *testing.T) {
	// Test complete engine with float64 data
	config := DefaultConfig()
	config.ProductionRate = 10 * time.Millisecond
	config.BatchSize = 5
	config.BatchTimeout = 50 * time.Millisecond
	config.MaxWorkers = 2

	// Simple time-based seeder
	seeder := NewTimeSeeder(1.0, 0.1, 0.0)

	// Simple sensor function
	sensorFunc := NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		return input*2.0 + 10.0
	})

	// Mock publisher
	publisher := &mockIntegrationPublisher[float64]{
		data: make([]SensorData[float64], 0),
	}

	// Create engine
	engine := NewEngine(config, seeder, sensorFunc, publisher)

	// Run for a short time
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Errorf("Engine start failed: %v", err)
	}

	// Verify data was generated
	if len(publisher.data) == 0 {
		t.Error("No data was published")
	}

	// Verify data structure
	for _, data := range publisher.data {
		if data.ID == "" {
			t.Error("Data ID should not be empty")
		}
		if data.Timestamp.IsZero() {
			t.Error("Data timestamp should not be zero")
		}
		if data.Quality == "" {
			t.Error("Data quality should not be empty")
		}
	}
}

func TestEngine_Integration_CustomStruct(t *testing.T) {
	// Test with custom data structure
	type CustomSensorData struct {
		Value     float64 `json:"value"`
		Unit      string  `json:"unit"`
		Location  string  `json:"location"`
		Timestamp int64   `json:"timestamp"`
	}

	config := DefaultConfig()
	config.ProductionRate = 20 * time.Millisecond
	config.BatchSize = 3
	config.BatchTimeout = 100 * time.Millisecond

	// Random seeder
	seeder := NewRandomSeeder(0.0, 1.0)

	// Custom sensor function
	sensorFunc := NewFunction(func(input float64, timestamp time.Time) CustomSensorData {
		return CustomSensorData{
			Value:     input * 100.0,
			Unit:      "percent",
			Location:  "sensor-001",
			Timestamp: timestamp.Unix(),
		}
	})

	// Mock publisher
	publisher := &mockIntegrationPublisher[CustomSensorData]{
		data: make([]SensorData[CustomSensorData], 0),
	}

	// Create engine
	engine := NewEngine(config, seeder, sensorFunc, publisher)

	// Run for a short time
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Errorf("Engine start failed: %v", err)
	}

	// Verify data was generated
	if len(publisher.data) == 0 {
		t.Error("No data was published")
	}

	// Verify custom data structure
	for _, data := range publisher.data {
		if data.Data.Unit != "percent" {
			t.Errorf("Expected unit 'percent', got '%s'", data.Data.Unit)
		}
		if data.Data.Location != "sensor-001" {
			t.Errorf("Expected location 'sensor-001', got '%s'", data.Data.Location)
		}
		if data.Data.Value < 0 || data.Data.Value > 100 {
			t.Errorf("Value %.2f out of range [0, 100]", data.Data.Value)
		}
	}
}

func TestEngine_Integration_Batching(t *testing.T) {
	// Test batching behavior
	config := DefaultConfig()
	config.ProductionRate = 5 * time.Millisecond
	config.BatchSize = 10
	config.BatchTimeout = 25 * time.Millisecond

	seeder := NewLinearSeeder(1.0, 0.0)

	sensorFunc := NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		return input
	})

	publisher := &mockIntegrationPublisher[float64]{
		data:       make([]SensorData[float64], 0),
		batchSizes: make([]int, 0),
	}

	engine := NewEngine(config, seeder, sensorFunc, publisher)

	// Run long enough for multiple batches
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Errorf("Engine start failed: %v", err)
	}

	// Verify batching occurred
	if len(publisher.batchSizes) == 0 {
		t.Error("No batches were published")
	}

	// Verify batch sizes are reasonable
	for _, batchSize := range publisher.batchSizes {
		if batchSize <= 0 || batchSize > config.BatchSize {
			t.Errorf("Batch size %d out of expected range [1, %d]", batchSize, config.BatchSize)
		}
	}
}

func TestEngine_Integration_QualitySimulation(t *testing.T) {
	// Test quality simulation
	config := DefaultConfig()
	config.ProductionRate = 10 * time.Millisecond
	config.BatchSize = 20

	seeder := NewTimeSeeder(1.0, 0.1, 0.0)

	sensorFunc := NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		return input
	})

	publisher := &mockIntegrationPublisher[float64]{
		data: make([]SensorData[float64], 0),
	}

	engine := NewEngine(config, seeder, sensorFunc, publisher)

	// Run long enough to see different quality levels
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Errorf("Engine start failed: %v", err)
	}

	// Verify quality levels are present
	qualityCounts := make(map[Quality]int)
	for _, data := range publisher.data {
		qualityCounts[data.Quality]++
	}

	// Should have at least some data with different quality levels
	if len(qualityCounts) == 0 {
		t.Error("No quality levels found")
	}

	// Most data should be OK quality
	okCount := qualityCounts[QualityOK]
	totalCount := len(publisher.data)
	if totalCount > 0 && float64(okCount)/float64(totalCount) < 0.5 {
		t.Errorf("Expected at least 50%% OK quality, got %.2f%%", float64(okCount)/float64(totalCount)*100)
	}
}

func TestEngine_Integration_ConcurrentAccess(t *testing.T) {
	// Test concurrent access to engine
	config := DefaultConfig()
	config.ProductionRate = 5 * time.Millisecond
	config.BatchSize = 5
	config.MaxWorkers = 3

	seeder := NewRandomSeeder(0.0, 1.0)

	sensorFunc := NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		// Simulate some processing time
		time.Sleep(1 * time.Millisecond)
		return input
	})

	publisher := &mockIntegrationPublisher[float64]{
		data: make([]SensorData[float64], 0),
	}

	engine := NewEngine(config, seeder, sensorFunc, publisher)

	// Run with concurrent access
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Errorf("Engine start failed: %v", err)
	}

	// Verify no race conditions occurred
	if len(publisher.data) == 0 {
		t.Error("No data was published during concurrent access")
	}

	// Verify data integrity
	for i, data := range publisher.data {
		if data.ID == "" {
			t.Errorf("Data %d has empty ID", i)
		}
		if data.Timestamp.IsZero() {
			t.Errorf("Data %d has zero timestamp", i)
		}
	}
}

func TestEngine_Integration_ErrorHandling(t *testing.T) {
	// Test error handling with failing publisher
	config := DefaultConfig()
	config.ProductionRate = 10 * time.Millisecond
	config.BatchSize = 3

	seeder := NewTimeSeeder(1.0, 0.1, 0.0)

	sensorFunc := NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		return input
	})

	// Publisher that always fails
	publisher := &failingMockPublisher[float64]{}

	engine := NewEngine(config, seeder, sensorFunc, publisher)

	// Run should handle errors gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// This should not panic, even though publisher fails
	err := engine.Start(ctx)
	// Engine should continue running despite publisher errors
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Unexpected engine error: %v", err)
	}
}

// Mock publisher for integration tests
type mockIntegrationPublisher[T any] struct {
	data       []SensorData[T]
	batchSizes []int
}

func (m *mockIntegrationPublisher[T]) Publish(ctx context.Context, data SensorData[T]) error {
	m.data = append(m.data, data)
	m.batchSizes = append(m.batchSizes, 1)
	return nil
}

func (m *mockIntegrationPublisher[T]) PublishBatch(ctx context.Context, data []SensorData[T]) error {
	m.data = append(m.data, data...)
	m.batchSizes = append(m.batchSizes, len(data))
	return nil
}

func (m *mockIntegrationPublisher[T]) Close() error {
	return nil
}

// Failing publisher for error handling tests
type failingMockPublisher[T any] struct{}

func (m *failingMockPublisher[T]) Publish(ctx context.Context, data SensorData[T]) error {
	return &mockError{"publish failed"}
}

func (m *failingMockPublisher[T]) PublishBatch(ctx context.Context, data []SensorData[T]) error {
	return &mockError{"batch publish failed"}
}

func (m *failingMockPublisher[T]) Close() error {
	return nil
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
