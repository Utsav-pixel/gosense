# 🚀 GoSense

A highly configurable, generic sensor data generation engine written in Go that can simulate any type of sensor data (medical, weather, industrial, etc.) with configurable production rates, batching, and multiple publishing options.

## 📖 **Documentation**

- **� [Documentation Hub](docs/)** - Complete documentation with guides and tutorials
- **�📝 Blog: [From Data Scarcity to Data Abundance](docs/BLOG_HUMAN.md)** - Learn how GoSense is revolutionizing software development
- **� [Seeder & Function Guide](docs/SEEDER_FUNCTION_GUIDE.md)** - Comprehensive guide on using seeders and functions
- **⚙️ [Advanced Configuration](docs/README_ADVANCED.md)** - Advanced setup and optimization techniques

## Features

- **Generic Type Support**: Works with any data type using Go generics
- **Flexible Seeders**: Multiple input generation strategies (time-based, random, linear, custom)
- **Configurable Functions**: Function-based sensor data generation
- **Multiple Publishers**: HTTP, Kafka, and gRPC support
- **Batch Processing**: Configurable batch sizes and timeouts
- **Production Rate Control**: Adjustable data generation frequency
- **Quality Simulation**: Realistic data quality variations
- **Concurrent Processing**: Multi-worker architecture for high throughput

## Architecture

The engine consists of several key components:

### 1. Generic Types (`internal/engine/types.go`)
- `SensorData[T]`: Generic container for sensor readings
- `Seeder`: Interface for input value generation
- `SensorFunction[T]`: Interface for data transformation
- `Publisher[T]`: Interface for data publishing
- `Engine[T]`: Main engine orchestrator

### 2. Seeders (`internal/engine/seeders.go`)
- `TimeSeeder`: Time-based oscillating values
- `RandomSeeder`: Random values within range
- `LinearSeeder`: Linearly increasing values
- `NormalSeeder`: Normal distribution values
- `CustomSeeder`: Custom generation functions

### 3. Sensor Functions (`internal/engine/functions.go`)
- `TemperatureSensorFunction`: Temperature data generation
- `HeartRateSensorFunction`: Heart rate simulation
- `BloodPressureSensorFunction`: Blood pressure readings
- `WeatherSensorFunction`: Weather data generation
- `CustomSensorFunction[T]`: Custom transformation functions

### 4. Publishers (`internal/publisher/`)
- `GenericHTTPPublisher[T]`: HTTP/REST API publishing
- `GenericKafkaPublisher[T]`: Apache Kafka publishing
- `GenericGRPCPublisher[T]`: gRPC streaming

## Quick Start

### Installation

```bash
git clone https://github.com/Utsav-pixel/go-sensor-engine.git
cd go-sensor-engine
go mod tidy
go build ./cmd/sensor-engine
```

### Running Examples

```bash
# Weather sensor with HTTP publisher (default)
./sensor-engine -type=weather -publisher=http -duration=30s

# Medical sensor with Kafka publisher
./sensor-engine -type=medical -publisher=kafka -brokers=localhost:9092 -topic=medical.data

# Industrial machinery sensor with gRPC publisher
./sensor-engine -type=machinery -publisher=grpc -grpc=localhost:50051

# Legacy pasture simulation
./sensor-engine -type=legacy
```

### Command Line Options

- `-type`: Sensor type (medical, weather, machinery, legacy)
- `-publisher`: Publisher type (http, kafka, grpc)
- `-duration`: How long to run the engine
- `-endpoint`: HTTP endpoint URL
- `-brokers`: Kafka broker addresses
- `-topic`: Kafka topic name
- `-grpc`: gRPC server address

## Usage Examples

### Medical Sensor Example

```go
// Configuration for medical sensors
config := engine.DefaultConfig()
config.ProductionRate = 1 * time.Second // Generate data every second
config.BatchSize = 10
config.BatchTimeout = 5 * time.Second

// Create a stress level seeder (0.0 to 1.0)
stressSeeder := engine.NewNormalSeeder(0.5, 0.2) // Patient variability

// Create medical sensor function with your own logic
medicalFunc := engine.NewFunction(func(input float64, timestamp time.Time) MedicalData {
    // Your business logic here
    heartRate := 70 + int(input*40) // Stress increases heart rate
    bloodPressure := BloodPressure{120 + int(input*20), 80 + int(input*15)}
    oxygenLevel := 95 + input*5 // Stress affects oxygen level
    temperature := 36.5 + input*2 // Stress affects temperature
    
    return MedicalData{
        HeartRate:     heartRate,
        BloodPressure: bloodPressure,
        OxygenLevel:   oxygenLevel,
        Temperature:   temperature,
    }
})

// Create publisher
httpPublisher := publisher.NewGenericHTTPPublisher[MedicalData]("https://api.medical.example.com/vitals")

// Create and start engine
medicalEngine := engine.NewEngine(config, stressSeeder, medicalFunc, httpPublisher)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := medicalEngine.Start(ctx); err != nil {
    log.Printf("Engine error: %v", err)
}
```

### Weather Sensor Example

```go
// High-throughput configuration
config := engine.HighThroughputConfig()
config.ProductionRate = 500 * time.Millisecond

// Weather pattern seeder
weatherSeeder := engine.NewTimeSeeder(1.0, 0.05, 0.0)

// Weather sensor function with your own logic
weatherFunc := engine.NewFunction(func(input float64, timestamp time.Time) WeatherData {
    // Your business logic here
    hour := float64(timestamp.Hour())
    dayOfYear := float64(timestamp.YearDay())
    
    // Temperature follows normal pattern with seasonal variation
    seasonalTemp := 15.0 + 10.0*math.Sin((dayOfYear/365.0)*2*math.Pi-math.Pi/2)
    dailyTemp := 5.0*math.Sin((hour/24.0)*2*math.Pi-math.Pi/2)
    temperature := seasonalTemp + dailyTemp + (input-0.5)*10.0
    
    // Humidity inversely related to temperature
    humidity := 70.0 - temperature + (rand.Float64()-0.5)*20.0
    if humidity < 20.0 { humidity = 20.0 } else if humidity > 95.0 { humidity = 95.0 }
    
    // Pressure varies with weather systems
    pressure := 1013.25 + (input-0.5)*50.0 + (rand.Float64()-0.5)*10.0
    
    // Wind speed
    windSpeed := math.Max(0, 10.0+input*30.0+(rand.Float64()-0.5)*5.0)
    
    return WeatherData{
        Temperature:  temperature,
        Humidity:    humidity,
        Pressure:    pressure,
        WindSpeed:   windSpeed,
    }
})

// Kafka publisher for high throughput
kafkaPublisher := publisher.NewGenericKafkaPublisher[WeatherData](
    []string{"localhost:9092"},
    "weather.data.v1",
)

// Create and start engine
weatherEngine := engine.NewEngine(config, weatherSeeder, weatherFunc, kafkaPublisher)
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
defer cancel()

if err := weatherEngine.Start(ctx); err != nil {
    log.Printf("Engine error: %v", err)
}
```

### Custom Sensor Example

```go
// Custom seeder that simulates market behavior
type MarketSeeder struct {
    cycle float64
}

func (m *MarketSeeder) Generate() float64 {
    m.cycle += 0.1
    baseValue := 0.5
    
    // Add market cycles
    cycle := math.Sin(m.cycle * 0.1) * 0.3
    
    // Add random market noise
    noise := (rand.Float64() - 0.5) * 0.2
    
    // Add trend component
    trend := math.Sin(m.cycle * 0.01) * 0.2
    
    result := baseValue + cycle + noise + trend
    
    // Keep within bounds
    if result < 0 { result = 0 } else if result > 1 { result = 1 }
    
    return result
}

// Custom sensor function with your own logic
customFunc := engine.NewFunction(func(input float64, timestamp time.Time) YourData {
    // Your business logic here
    value := input * 100.0
    status := "normal"
    if value > 80 { status = "high" } else if value > 60 { status = "medium" }
    
    return YourData{
        Value:    value,
        Status:   status,
        Location: "sensor-001",
    }
})

// Create publisher
consolePublisher := NewConsolePublisher[YourData]()

// Create and start engine
customEngine := engine.NewEngine(config, &MarketSeeder{cycle: 0}, customFunc, consolePublisher)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := customEngine.Start(ctx); err != nil {
    log.Printf("Engine error: %v", err)
}
```

## Configuration Presets

### DefaultConfig
- Production Rate: 100ms
- Batch Size: 100
- Batch Timeout: 500ms
- Max Workers: 3

### HighThroughputConfig
- Production Rate: 10ms
- Batch Size: 1000
- Batch Timeout: 100ms
- Max Workers: 10

### LowLatencyConfig
- Production Rate: 50ms
- Batch Size: 10
- Batch Timeout: 25ms
- Max Workers: 5

## Data Quality

The engine simulates realistic data quality variations:

- **OK** (89%): Normal quality data
- **NOISY** (5%): Data with some noise
- **PARTIAL** (2%): Partially complete data
- **CORRUPT** (1%): Corrupted data

## Extending the Engine

### Adding Custom Seeders

```go
type MyCustomSeeder struct {
    // your fields
}

func (m *MyCustomSeeder) Generate() float64 {
    // your implementation
}
```

### Adding Custom Sensor Functions

```go
type MyCustomSensorFunction struct {
    // your fields
}

func (m *MyCustomSensorFunction) Generate(input float64, timestamp time.Time) MyDataType {
    // your implementation
}
```

### Adding Custom Publishers

```go
type MyCustomPublisher[T any] struct {
    // your fields
}

func (m *MyCustomPublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
    // your implementation
}

func (m *MyCustomPublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
    // your implementation
}

func (m *MyCustomPublisher[T]) Close() error {
    // your implementation
}
```

## Performance Considerations

- Use `HighThroughputConfig` for high-volume data generation
- Use `LowLatencyConfig` for real-time applications
- Adjust batch sizes based on your downstream system capacity
- Consider using Kafka for high-throughput scenarios
- Use gRPC for low-latency real-time monitoring

## Dependencies

- Go 1.24+
- google.golang.org/grpc v1.65.0
- google.golang.org/protobuf v1.34.2
- github.com/segmentio/kafka-go v0.4.50

## License

This project is licensed under the MIT License - see the LICENSE file for details.
