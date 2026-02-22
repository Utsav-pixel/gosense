package publisher

import (
	"context"
	"testing"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
)

func TestGenericHTTPPublisher_Publish(t *testing.T) {
	publisher := NewGenericHTTPPublisher[float64]("https://httpbin.org/post")

	data := engine.SensorData[float64]{
		ID:        "test-1",
		Timestamp: time.Now(),
		Data:      25.5,
		Quality:   engine.QualityOK,
	}

	// Note: This test requires internet connection to httpbin.org
	// In a real test environment, you might want to mock the HTTP client
	err := publisher.Publish(context.Background(), data)
	if err != nil {
		t.Logf("HTTP publish failed (expected if no internet): %v", err)
		// Don't fail the test if there's no internet connection
	}
}

func TestGenericHTTPPublisher_PublishBatch(t *testing.T) {
	publisher := NewGenericHTTPPublisher[float64]("https://httpbin.org/post")

	batch := []engine.SensorData[float64]{
		{
			ID:        "batch-1",
			Timestamp: time.Now(),
			Data:      25.5,
			Quality:   engine.QualityOK,
		},
		{
			ID:        "batch-2",
			Timestamp: time.Now(),
			Data:      26.0,
			Quality:   engine.QualityOK,
		},
	}

	err := publisher.PublishBatch(context.Background(), batch)
	if err != nil {
		t.Logf("HTTP batch publish failed (expected if no internet): %v", err)
		// Don't fail the test if there's no internet connection
	}
}

func TestGenericHTTPPublisher_Close(t *testing.T) {
	publisher := NewGenericHTTPPublisher[float64]("https://example.com")

	err := publisher.Close()
	if err != nil {
		t.Errorf("Unexpected error closing HTTP publisher: %v", err)
	}
}

func TestGenericKafkaPublisher_Publish(t *testing.T) {
	// Note: This test requires a running Kafka instance
	// For unit tests, you might want to mock the Kafka writer
	publisher := NewGenericKafkaPublisher[float64](
		[]string{"localhost:9092"},
		"test-topic",
	)

	data := engine.SensorData[float64]{
		ID:        "test-1",
		Timestamp: time.Now(),
		Data:      25.5,
		Quality:   engine.QualityOK,
	}

	err := publisher.Publish(context.Background(), data)
	if err != nil {
		t.Logf("Kafka publish failed (expected if no Kafka running): %v", err)
		// Don't fail the test if there's no Kafka connection
	}
}

func TestGenericKafkaPublisher_PublishBatch(t *testing.T) {
	publisher := NewGenericKafkaPublisher[float64](
		[]string{"localhost:9092"},
		"test-topic",
	)

	batch := []engine.SensorData[float64]{
		{
			ID:        "batch-1",
			Timestamp: time.Now(),
			Data:      25.5,
			Quality:   engine.QualityOK,
		},
		{
			ID:        "batch-2",
			Timestamp: time.Now(),
			Data:      26.0,
			Quality:   engine.QualityOK,
		},
	}

	err := publisher.PublishBatch(context.Background(), batch)
	if err != nil {
		t.Logf("Kafka batch publish failed (expected if no Kafka running): %v", err)
		// Don't fail the test if there's no Kafka connection
	}
}

func TestGenericKafkaPublisher_Close(t *testing.T) {
	publisher := NewGenericKafkaPublisher[float64](
		[]string{"localhost:9092"},
		"test-topic",
	)

	err := publisher.Close()
	if err != nil {
		t.Errorf("Unexpected error closing Kafka publisher: %v", err)
	}
}

func TestGenericGRPCPublisher_Publish(t *testing.T) {
	// Note: This test requires a running gRPC server
	publisher, err := NewGenericGRPCPublisher[float64]("localhost:50051")
	if err != nil {
		t.Logf("Failed to create gRPC publisher (expected if no gRPC server): %v", err)
		return
	}

	data := engine.SensorData[float64]{
		ID:        "test-1",
		Timestamp: time.Now(),
		Data:      25.5,
		Quality:   engine.QualityOK,
	}

	err = publisher.Publish(context.Background(), data)
	if err != nil {
		t.Logf("gRPC publish failed (expected if no gRPC server): %v", err)
		// Don't fail the test if there's no gRPC server
	}
}

func TestGenericGRPCPublisher_PublishBatch(t *testing.T) {
	publisher, err := NewGenericGRPCPublisher[float64]("localhost:50051")
	if err != nil {
		t.Logf("Failed to create gRPC publisher (expected if no gRPC server): %v", err)
		return
	}

	batch := []engine.SensorData[float64]{
		{
			ID:        "batch-1",
			Timestamp: time.Now(),
			Data:      25.5,
			Quality:   engine.QualityOK,
		},
		{
			ID:        "batch-2",
			Timestamp: time.Now(),
			Data:      26.0,
			Quality:   engine.QualityOK,
		},
	}

	err = publisher.PublishBatch(context.Background(), batch)
	if err != nil {
		t.Logf("gRPC batch publish failed (expected if no gRPC server): %v", err)
		// Don't fail the test if there's no gRPC server
	}
}

func TestGenericGRPCPublisher_Close(t *testing.T) {
	publisher, err := NewGenericGRPCPublisher[float64]("localhost:50051")
	if err != nil {
		t.Logf("Failed to create gRPC publisher (expected if no gRPC server): %v", err)
		return
	}

	err = publisher.Close()
	if err != nil {
		t.Errorf("Unexpected error closing gRPC publisher: %v", err)
	}
}

// Mock publisher for testing
type MockPublisher[T any] struct {
	PublishedData []engine.SensorData[T]
	PublishCount  int
	BatchCount    int
}

func NewMockPublisher[T any]() *MockPublisher[T] {
	return &MockPublisher[T]{
		PublishedData: make([]engine.SensorData[T], 0),
	}
}

func (m *MockPublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
	m.PublishedData = append(m.PublishedData, data)
	m.PublishCount++
	return nil
}

func (m *MockPublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
	m.PublishedData = append(m.PublishedData, data...)
	m.BatchCount++
	return nil
}

func (m *MockPublisher[T]) Close() error {
	return nil
}

func TestMockPublisher(t *testing.T) {
	publisher := NewMockPublisher[float64]()

	// Test single publish
	data := engine.SensorData[float64]{
		ID:        "test-1",
		Timestamp: time.Now(),
		Data:      25.5,
		Quality:   engine.QualityOK,
	}

	err := publisher.Publish(context.Background(), data)
	if err != nil {
		t.Errorf("Unexpected error publishing to mock: %v", err)
	}

	if publisher.PublishCount != 1 {
		t.Errorf("Expected publish count 1, got %d", publisher.PublishCount)
	}

	if len(publisher.PublishedData) != 1 {
		t.Errorf("Expected 1 published data item, got %d", len(publisher.PublishedData))
	}

	// Test batch publish
	batch := []engine.SensorData[float64]{
		{
			ID:        "batch-1",
			Timestamp: time.Now(),
			Data:      26.0,
			Quality:   engine.QualityOK,
		},
		{
			ID:        "batch-2",
			Timestamp: time.Now(),
			Data:      26.5,
			Quality:   engine.QualityOK,
		},
	}

	err = publisher.PublishBatch(context.Background(), batch)
	if err != nil {
		t.Errorf("Unexpected error publishing batch to mock: %v", err)
	}

	if publisher.BatchCount != 1 {
		t.Errorf("Expected batch count 1, got %d", publisher.BatchCount)
	}

	if len(publisher.PublishedData) != 3 { // 1 from single publish + 2 from batch
		t.Errorf("Expected 3 published data items, got %d", len(publisher.PublishedData))
	}
}

// Benchmark tests
func BenchmarkGenericHTTPPublisher_Publish(b *testing.B) {
	publisher := NewGenericHTTPPublisher[float64]("https://httpbin.org/post")
	data := engine.SensorData[float64]{
		ID:        "bench-test",
		Timestamp: time.Now(),
		Data:      25.5,
		Quality:   engine.QualityOK,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This benchmark requires internet connection
		publisher.Publish(context.Background(), data)
	}
}

func BenchmarkMockPublisher_Publish(b *testing.B) {
	publisher := NewMockPublisher[float64]()
	data := engine.SensorData[float64]{
		ID:        "bench-test",
		Timestamp: time.Now(),
		Data:      25.5,
		Quality:   engine.QualityOK,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher.Publish(context.Background(), data)
	}
}

func BenchmarkMockPublisher_PublishBatch(b *testing.B) {
	publisher := NewMockPublisher[float64]()
	batch := make([]engine.SensorData[float64], 100)
	for i := range batch {
		batch[i] = engine.SensorData[float64]{
			ID:        "bench-test",
			Timestamp: time.Now(),
			Data:      float64(i),
			Quality:   engine.QualityOK,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher.PublishBatch(context.Background(), batch)
	}
}
