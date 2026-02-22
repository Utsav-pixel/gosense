package engine

import (
	"context"
	"testing"
	"time"
)

// MockPublisher for testing
type MockPublisher[T any] struct {
	published     []SensorData[T]
	batches       [][]SensorData[T]
	publishCalled int
	batchCalled   int
	closed        bool
}

func NewMockPublisher[T any]() *MockPublisher[T] {
	return &MockPublisher[T]{
		published: make([]SensorData[T], 0),
		batches:   make([][]SensorData[T], 0),
	}
}

func (m *MockPublisher[T]) Publish(ctx context.Context, data SensorData[T]) error {
	m.published = append(m.published, data)
	m.publishCalled++
	return nil
}

func (m *MockPublisher[T]) PublishBatch(ctx context.Context, data []SensorData[T]) error {
	m.batches = append(m.batches, data)
	m.batchCalled++
	return nil
}

func (m *MockPublisher[T]) Close() error {
	m.closed = true
	return nil
}

func (m *MockPublisher[T]) GetPublishedCount() int {
	return m.publishCalled
}

func (m *MockPublisher[T]) GetBatchCount() int {
	return m.batchCalled
}

func (m *MockPublisher[T]) GetTotalDataPoints() int {
	total := len(m.published)
	for _, batch := range m.batches {
		total += len(batch)
	}
	return total
}

func (m *MockPublisher[T]) IsClosed() bool {
	return m.closed
}

// TestSeeder for testing
type TestSeeder struct {
	values []float64
	index  int
}

func NewTestSeeder(values []float64) *TestSeeder {
	return &TestSeeder{
		values: values,
		index:  0,
	}
}

func (t *TestSeeder) Generate() float64 {
	if t.index >= len(t.values) {
		t.index = 0 // Reset for cyclic behavior
	}
	val := t.values[t.index]
	t.index++
	return val
}

// TestSensorFunction for testing
type TestSensorFunction struct {
	multiplier float64
}

func NewTestSensorFunction(multiplier float64) *TestSensorFunction {
	return &TestSensorFunction{multiplier: multiplier}
}

func (t *TestSensorFunction) Generate(input float64, timestamp time.Time) float64 {
	return input * t.multiplier
}

func TestEngine_BasicFunctionality(t *testing.T) {
	// Setup
	config := Config{
		ProductionRate: 10 * time.Millisecond,
		BatchSize:      2,
		BatchTimeout:   50 * time.Millisecond,
		MaxWorkers:     1,
	}

	seeder := NewTestSeeder([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
	function := NewTestSensorFunction(2.0)
	publisher := NewMockPublisher[float64]()

	engine := NewEngine(config, seeder, function, publisher)

	// Run for short duration
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Fatalf("Engine start failed: %v", err)
	}

	// Verify results
	if publisher.GetTotalDataPoints() == 0 {
		t.Error("No data was published")
	}

	if !publisher.IsClosed() {
		t.Error("Publisher was not closed")
	}

	t.Logf("Published %d data points in %d batches",
		publisher.GetTotalDataPoints(), publisher.GetBatchCount())
}

func TestEngine_BatchProcessing(t *testing.T) {
	// Setup with specific batch configuration
	config := Config{
		ProductionRate: 5 * time.Millisecond,
		BatchSize:      3,
		BatchTimeout:   20 * time.Millisecond,
		MaxWorkers:     1,
	}

	seeder := NewTestSeeder([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0})
	function := NewTestSensorFunction(1.5)
	publisher := NewMockPublisher[float64]()

	engine := NewEngine(config, seeder, function, publisher)

	// Run for enough time to generate multiple batches
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Fatalf("Engine start failed: %v", err)
	}

	// Verify batch processing
	if publisher.GetBatchCount() == 0 {
		t.Error("No batches were published")
	}

	// Check that batches contain expected number of items
	for i, batch := range publisher.batches {
		if len(batch) > config.BatchSize {
			t.Errorf("Batch %d has %d items, expected max %d", i, len(batch), config.BatchSize)
		}
	}

	t.Logf("Processed %d batches with total %d data points",
		publisher.GetBatchCount(), publisher.GetTotalDataPoints())
}

func TestEngine_QualityGeneration(t *testing.T) {
	config := DefaultConfig()
	config.ProductionRate = 5 * time.Millisecond
	config.BatchSize = 1

	seeder := NewTestSeeder([]float64{1.0})
	function := NewTestSensorFunction(1.0)
	publisher := NewMockPublisher[float64]()

	engine := NewEngine(config, seeder, function, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		t.Fatalf("Engine start failed: %v", err)
	}

	// Verify that quality is being set
	totalData := publisher.GetTotalDataPoints()
	if totalData == 0 {
		t.Error("No data was published")
	}

	// Check that all data has quality set
	for _, batch := range publisher.batches {
		for _, data := range batch {
			if data.Quality == "" {
				t.Error("Data quality was not set")
			}
		}
	}

	t.Logf("Generated %d data points with quality", totalData)
}

func TestEngine_ContextCancellation(t *testing.T) {
	config := DefaultConfig()
	seeder := NewTestSeeder([]float64{1.0, 2.0, 3.0})
	function := NewTestSensorFunction(1.0)
	publisher := NewMockPublisher[float64]()

	engine := NewEngine(config, seeder, function, publisher)

	// Create a context that cancels quickly
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short time
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := engine.Start(ctx)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Engine start failed: %v", err)
	}

	// Should stop quickly after cancellation
	if duration > 100*time.Millisecond {
		t.Errorf("Engine took too long to stop: %v", duration)
	}

	if !publisher.IsClosed() {
		t.Error("Publisher was not closed")
	}

	t.Logf("Engine stopped in %v after context cancellation", duration)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.ProductionRate == 0 {
		t.Error("Production rate not set")
	}
	if config.BatchSize == 0 {
		t.Error("Batch size not set")
	}
	if config.BatchTimeout == 0 {
		t.Error("Batch timeout not set")
	}
	if config.MaxWorkers == 0 {
		t.Error("Max workers not set")
	}
}

func TestHighThroughputConfig(t *testing.T) {
	config := HighThroughputConfig()
	defaultConfig := DefaultConfig()

	if config.ProductionRate >= defaultConfig.ProductionRate {
		t.Error("High throughput config should have faster production rate")
	}
	if config.BatchSize <= defaultConfig.BatchSize {
		t.Error("High throughput config should have larger batch size")
	}
	if config.BatchTimeout >= defaultConfig.BatchTimeout {
		t.Error("High throughput config should have shorter batch timeout")
	}
	if config.MaxWorkers <= defaultConfig.MaxWorkers {
		t.Error("High throughput config should have more workers")
	}
}

func TestLowLatencyConfig(t *testing.T) {
	config := LowLatencyConfig()
	defaultConfig := DefaultConfig()

	if config.BatchSize >= defaultConfig.BatchSize {
		t.Error("Low latency config should have smaller batch size")
	}
	if config.BatchTimeout >= defaultConfig.BatchTimeout {
		t.Error("Low latency config should have shorter batch timeout")
	}
}

// Benchmark tests
func BenchmarkEngine_DataGeneration(b *testing.B) {
	config := Config{
		ProductionRate: 1 * time.Millisecond,
		BatchSize:      100,
		BatchTimeout:   10 * time.Millisecond,
		MaxWorkers:     4,
	}

	seeder := NewTestSeeder([]float64{1.0})
	function := NewTestSensorFunction(1.0)
	publisher := NewMockPublisher[float64]()

	engine := NewEngine(config, seeder, function, publisher)

	b.ResetTimer()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := engine.Start(ctx)
	if err != nil {
		b.Fatalf("Engine start failed: %v", err)
	}

	b.StopTimer()
	dataPoints := publisher.GetTotalDataPoints()
	b.Logf("Generated %d data points in 1 second", dataPoints)
	b.ReportMetric(float64(dataPoints), "data_points/sec")
}
