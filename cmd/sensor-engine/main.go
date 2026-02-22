package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/examples"
	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
)

func main() {
	var (
		sensorType = flag.String("type", "", "Sensor example type: temperature, iot, industrial, weather, financial, config")
		config     = flag.String("config", "", "JSON configuration file path")
		duration   = flag.Duration("duration", 10*time.Second, "How long to run the sensor engine")
		help       = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *sensorType == "" && *config == "" {
		fmt.Println("Error: Please specify either -type or -config")
		showHelp()
		os.Exit(1)
	}

	if *config != "" {
		runFromConfig(*config, *duration)
		return
	}

	switch *sensorType {
	case "temperature":
		log.Println("🌡️  Starting Temperature Sensor Example...")
		examples.TemperatureSensorExample()
	case "iot":
		log.Println("📱 Starting IoT Device Example...")
		examples.IoTDeviceExample()
	case "industrial":
		log.Println("🏭 Starting Industrial Sensor Example...")
		examples.IndustrialSensorExample()
	case "weather":
		log.Println("🌤️  Starting Weather Station Example...")
		examples.WeatherStationExample()
	case "financial":
		log.Println("💰 Starting Financial Metrics Example...")
		examples.CustomSeederExample()
	default:
		fmt.Printf("Error: Unknown sensor type '%s'\n", *sensorType)
		showHelp()
		os.Exit(1)
	}
}

func runFromConfig(configPath string, duration time.Duration) {
	log.Printf("🚀 Starting sensor engine from config: %s", configPath)

	// Create a simple function for demonstration
	sensorFunc := engine.NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
		return input * 100.0
	})

	// Create engine from config
	testEngine, err := engine.CreateEngineFromConfig(configPath, sensorFunc, examples.NewConsolePublisher[float64]())
	if err != nil {
		log.Fatalf("Failed to create engine from config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}

	log.Println("✅ Sensor engine completed successfully")
}

func showHelp() {
	fmt.Println(`
🎯 Generic Sensor Engine - Real-World Examples

USAGE:
  sensor-engine -type=<example_type> [options]
  sensor-engine -config=<config_file> [options]

EXAMPLE TYPES:
  temperature    🌡️  Temperature sensor with time-based seeder showing daily cycles
  iot            📱 IoT device with random seeder simulating battery and signal
  industrial     🏭 Industrial sensor with linear seeder showing machine wear
  weather        🌤️  Weather station with normal distribution seeder
  financial      💰 Financial metrics with custom market behavior seeder

OPTIONS:
  -type <type>        Example type to run (see list above)
  -config <file>      JSON configuration file to use
  -publisher <type>    Publisher type (console, http, kafka, grpc)
  -duration <time>     How long to run (default: 10s)
  -help               Show this help message

SEEDER + FUNCTION INTEGRATION EXAMPLES:

1. Temperature Sensor (Time-based Seeder):
   - Seeder: NewTimeSeeder(1.0, 0.1, 20.0)
   - Function: Uses seeder input + time for daily temperature cycles
   - Shows: How environmental factors change over time

2. IoT Device (Random Seeder):
   - Seeder: NewRandomSeeder(0.0, 1.0)  
   - Function: Maps random input to device metrics (battery, signal, etc.)
   - Shows: How random events trigger sensor readings

3. Industrial Sensor (Linear Seeder):
   - Seeder: NewLinearSeeder(0.01, 0.1)
   - Function: Simulates machine degradation over time
   - Shows: Progressive changes and wear patterns

4. Weather Station (Normal Seeder):
   - Seeder: NewNormalSeeder(0.5, 0.2)
   - Function: Natural weather variations with realistic patterns
   - Shows: Statistical distributions in natural phenomena

5. Financial Metrics (Custom Seeder):
   - Seeder: Custom MarketSeeder with cycles and trends
   - Function: Market sentiment to financial data transformation
   - Shows: Complex custom seeder + function combinations

JSON CONFIGURATION:
  Use -config with any JSON file in configs/ directory:
  - configs/temperature-sensor.json
  - configs/medical-sensor.json  
  - configs/industrial-sensor.json

EXAMPLES:
  # Run temperature sensor for 30 seconds
  sensor-engine -type=temperature -duration=30s

  # Run IoT device with HTTP publisher
  sensor-engine -type=iot -publisher=http -duration=1m

  # Run from JSON configuration
  sensor-engine -config=configs/temperature-sensor.json -duration=2m

  # Run financial metrics example
  sensor-engine -type=financial -duration=45s
`)
}
