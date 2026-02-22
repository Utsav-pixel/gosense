package engine

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeSeeder(t *testing.T) {
	seeder := NewTimeSeeder(1.0, 1.0, 0.5) // Higher frequency for faster changes

	value1 := seeder.Generate()
	time.Sleep(10 * time.Millisecond) // Shorter sleep is fine with higher frequency
	value2 := seeder.Generate()

	if value1 == value2 {
		t.Error("TimeSeeder should generate different values over time")
	}

	// Test that values are within expected range
	if value1 < -0.5 || value1 > 1.5 { // amplitude=1, offset=0.5, so range is [-0.5, 1.5]
		t.Errorf("Value %f outside expected range [-0.5, 1.5]", value1)
	}
}

func TestRandomSeeder(t *testing.T) {
	min, max := 10.0, 20.0
	seeder := NewRandomSeeder(min, max)

	// Generate multiple values
	for i := 0; i < 100; i++ {
		value := seeder.Generate()
		if value < min || value > max {
			t.Errorf("Value %f outside range [%f, %f]", value, min, max)
		}
	}
}

func TestLinearSeeder(t *testing.T) {
	slope, offset := 2.0, 10.0
	seeder := NewLinearSeeder(slope, offset)

	value1 := seeder.Generate()
	time.Sleep(10 * time.Millisecond)
	value2 := seeder.Generate()

	// Should be increasing
	if value2 <= value1 {
		t.Error("LinearSeeder should generate increasing values")
	}

	// Test initial value is close to offset
	if value1 < offset-1.0 || value1 > offset+1.0 {
		t.Errorf("Initial value %f too far from offset %f", value1, offset)
	}
}

func TestNormalSeeder(t *testing.T) {
	mean, stdDev := 50.0, 10.0
	seeder := NewNormalSeeder(mean, stdDev)

	// Generate multiple values and check distribution
	sum := 0.0
	count := 1000
	for i := 0; i < count; i++ {
		value := seeder.Generate()
		sum += value
	}

	avg := sum / float64(count)

	// Average should be close to mean (within 5% for 1000 samples)
	if avg < mean*0.95 || avg > mean*1.05 {
		t.Errorf("Average %f too far from mean %f", avg, mean)
	}
}

func TestCustomSeeder(t *testing.T) {
	calls := 0
	seeder := NewCustomSeeder(func() float64 {
		calls++
		return float64(calls) * 2.0
	})

	value1 := seeder.Generate()
	value2 := seeder.Generate()

	if value1 != 2.0 {
		t.Errorf("Expected first call to return 2.0, got %f", value1)
	}

	if value2 != 4.0 {
		t.Errorf("Expected second call to return 4.0, got %f", value2)
	}
}

func TestBasicSensorFunction(t *testing.T) {
	// Test with string output
	function := NewBasicSensorFunction(func(input float64, timestamp time.Time) string {
		return fmt.Sprintf("value_%f_at_%v", input, timestamp.Unix())
	})

	result := function.Generate(1.5, time.Unix(1234567890, 0))
	expected := "value_1.500000_at_1234567890"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestCustomSensorFunction(t *testing.T) {
	type TestData struct {
		Value     float64 `json:"value"`
		Timestamp int64   `json:"timestamp"`
		Processed bool    `json:"processed"`
	}

	function := NewFunction(func(input float64, timestamp time.Time) TestData {
		return TestData{
			Value:     input * 2.0,
			Timestamp: timestamp.Unix(),
			Processed: true,
		}
	})

	ts := time.Unix(1234567890, 0)
	result := function.Generate(2.5, ts)

	if result.Value != 5.0 {
		t.Errorf("Expected Value 5.0, got %f", result.Value)
	}

	if result.Timestamp != 1234567890 {
		t.Errorf("Expected Timestamp 1234567890, got %d", result.Timestamp)
	}

	if !result.Processed {
		t.Error("Expected Processed to be true")
	}
}

func TestLambdaSensorFunction(t *testing.T) {
	// Test with inline lambda
	function := NewLambdaSensorFunction(func(input float64, timestamp time.Time) int {
		return int(input * 100)
	})

	result := function.Generate(1.234, time.Now())

	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}
}

// Benchmark seeders
func BenchmarkTimeSeeder(b *testing.B) {
	seeder := NewTimeSeeder(1.0, 0.1, 0.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seeder.Generate()
	}
}

func BenchmarkRandomSeeder(b *testing.B) {
	seeder := NewRandomSeeder(0.0, 100.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seeder.Generate()
	}
}

func BenchmarkNormalSeeder(b *testing.B) {
	seeder := NewNormalSeeder(50.0, 10.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seeder.Generate()
	}
}

func BenchmarkBasicSensorFunction(b *testing.B) {
	function := NewBasicSensorFunction(func(input float64, timestamp time.Time) float64 {
		return input * 2.0
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		function.Generate(1.0, time.Now())
	}
}
