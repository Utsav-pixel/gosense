package examples

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
)

// ConsolePublisher for testing and demonstration
type ConsolePublisher[T any] struct{}

func NewConsolePublisher[T any]() *ConsolePublisher[T] {
	return &ConsolePublisher[T]{}
}

func (p *ConsolePublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
	fmt.Printf("📊 [%s] %+v\n", data.Quality, data.Data)
	return nil
}

func (p *ConsolePublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
	fmt.Printf("📦 Batch of %d items:\n", len(data))
	for i, item := range data {
		fmt.Printf("  [%d] [%s] %+v\n", i, item.Quality, item.Data)
	}
	return nil
}

func (p *ConsolePublisher[T]) Close() error {
	fmt.Println("🔚 Console publisher closed")
	return nil
}

// Example 1: Temperature Sensor with Time-based Seeder
// Shows how environmental factors change over time
func TemperatureSensorExample() {
	type TemperatureReading struct {
		Celsius    float64 `json:"celsius"`
		Fahrenheit float64 `json:"fahrenheit"`
		Humidity   float64 `json:"humidity_percent"`
		Location   string  `json:"location"`
	}

	// Time-based seeder generates values that change over time
	// This simulates daily temperature cycles
	seeder := engine.NewTimeSeeder(
		1.0,  // amplitude - temperature variation range
		0.1,  // frequency - how fast temperature changes
		20.0, // offset - base temperature
	)

	// User-defined function that uses seeder input to generate realistic temperature data
	sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) TemperatureReading {
		// Input from seeder represents environmental factor (0-1) affecting temperature
		// Higher input = hotter environment

		// Base temperature varies with seeder input
		baseTemp := input

		// Add diurnal pattern (daily temperature cycle)
		hour := float64(timestamp.Hour()) + float64(timestamp.Minute())/60.0
		radian := (hour / 24.0) * 2 * math.Pi
		diurnal := 5.0 * math.Sin(radian-math.Pi/2) // Peak at 2 PM

		// Add random noise for realism
		noise := (rand.Float64() - 0.5) * 1.0

		celsius := baseTemp + diurnal + noise
		fahrenheit := celsius*9/5 + 32

		// Humidity inversely related to temperature
		humidity := 70.0 - celsius
		if humidity < 30.0 {
			humidity = 30.0
		} else if humidity > 90.0 {
			humidity = 90.0
		}

		return TemperatureReading{
			Celsius:    celsius,
			Fahrenheit: fahrenheit,
			Humidity:   humidity,
			Location:   "Server Room A",
		}
	})

	publisher := NewConsolePublisher[TemperatureReading]()

	config := engine.DefaultConfig()
	config.ProductionRate = 1 * time.Second
	config.BatchSize = 3

	testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("🌡️  Starting Temperature Sensor Example...")
	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}
}

// Example 2: IoT Device with Random Seeder
// Shows how random events can trigger sensor readings
func IoTDeviceExample() {
	type IoTReading struct {
		DeviceID    string  `json:"device_id"`
		Battery     float64 `json:"battery_percent"`
		Signal      int     `json:"signal_strength_dbm"`
		Temperature float64 `json:"temperature_celsius"`
		Status      string  `json:"status"`
		LastSeen    int64   `json:"last_seen_unix"`
	}

	// Random seeder simulates random device states
	// Each call generates a random value between 0 and 1
	seeder := engine.NewRandomSeeder(0.0, 1.0)

	// User-defined function that uses random input to simulate IoT device behavior
	sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) IoTReading {
		// Input from random seeder represents device stress/activity level
		// Higher input = more active device = higher battery drain

		// Battery level decreases with activity
		battery := 100.0 - (input * 30.0) // 70-100% range

		// Signal strength varies with activity (movement affects signal)
		signal := -30 - int(input*40) // -30 to -70 dBm

		// Temperature increases with device activity
		temperature := 25.0 + (input * 15.0) // 25-40°C

		// Device status based on battery level
		var status string
		switch {
		case battery > 80:
			status = "excellent"
		case battery > 50:
			status = "good"
		case battery > 20:
			status = "low"
		default:
			status = "critical"
		}

		return IoTReading{
			DeviceID:    fmt.Sprintf("iot-%04d", int(input*9999)),
			Battery:     battery,
			Signal:      signal,
			Temperature: temperature,
			Status:      status,
			LastSeen:    timestamp.Unix(),
		}
	})

	publisher := NewConsolePublisher[IoTReading]()

	config := engine.DefaultConfig()
	config.ProductionRate = 500 * time.Millisecond
	config.BatchSize = 5

	testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	log.Println("📱 Starting IoT Device Example...")
	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}
}

// Example 3: Industrial Sensor with Linear Seeder
// Shows how progressive changes can be simulated
func IndustrialSensorExample() {
	type MachineMetrics struct {
		MachineID   string  `json:"machine_id"`
		Vibration   float64 `json:"vibration_mm_s"`
		Pressure    float64 `json:"pressure_bar"`
		RPM         int     `json:"rpm"`
		Temperature float64 `json:"temperature_celsius"`
		Efficiency  float64 `json:"efficiency_percent"`
		Status      string  `json:"status"`
	}

	// Linear seeder simulates gradual machine wear over time
	// Starts at 0.1 and increases by 0.01 each generation
	seeder := engine.NewLinearSeeder(0.01, 0.1)

	// User-defined function that simulates machine degradation
	sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) MachineMetrics {
		// Input from linear seeder represents machine wear factor (0.1 to 1.0+)
		// Higher input = more wear = worse performance

		// Vibration increases with wear
		vibration := input * 8.0 // 0.8 to 8+ mm/s

		// System pressure becomes unstable with wear
		pressure := 2.5 + (input * 1.5) + (rand.Float64()-0.5)*0.5

		// RPM decreases as machine wears out
		baseRPM := 1800
		rpmReduction := int(input * 400) // Up to 400 RPM reduction
		rpm := baseRPM - rpmReduction

		// Temperature increases with friction from wear
		temperature := 25.0 + (input * 30.0)

		// Efficiency decreases with wear
		efficiency := 100.0 - (input * 40.0) // 60% to 100%
		if efficiency < 0 {
			efficiency = 0
		}

		// Status based on wear level
		var status string
		switch {
		case input > 0.8:
			status = "critical_maintenance"
		case input > 0.6:
			status = "warning"
		case input > 0.3:
			status = "monitor"
		default:
			status = "normal"
		}

		return MachineMetrics{
			MachineID:   fmt.Sprintf("CNC-%03d", int(input*999)),
			Vibration:   vibration,
			Pressure:    pressure,
			RPM:         rpm,
			Temperature: temperature,
			Efficiency:  efficiency,
			Status:      status,
		}
	})

	publisher := NewConsolePublisher[MachineMetrics]()

	config := engine.DefaultConfig()
	config.ProductionRate = 2 * time.Second
	config.BatchSize = 2

	testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	log.Println("🏭 Starting Industrial Sensor Example...")
	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}
}

// Example 4: Weather Station with Normal Distribution Seeder
// Shows how realistic statistical patterns can be simulated
func WeatherStationExample() {
	type WeatherData struct {
		StationID     string  `json:"station_id"`
		Temperature   float64 `json:"temperature_celsius"`
		Humidity      float64 `json:"humidity_percent"`
		Pressure      float64 `json:"pressure_hpa"`
		WindSpeed     float64 `json:"wind_speed_kmh"`
		WindDirection int     `json:"wind_direction_degrees"`
		Conditions    string  `json:"conditions"`
		Timestamp     int64   `json:"timestamp_unix"`
	}

	// Normal seeder generates values following normal distribution
	// Mean=0.5, StdDev=0.2 - simulates natural weather variations
	seeder := engine.NewNormalSeeder(0.5, 0.2)

	// User-defined function that simulates realistic weather patterns
	sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) WeatherData {
		// Input from normal seeder represents weather variability factor
		// Most values cluster around 0.5 with some outliers

		// Temperature follows normal pattern with seasonal variation
		hour := float64(timestamp.Hour())
		dayOfYear := float64(timestamp.YearDay())

		// Base temperature with seasonal variation
		seasonalTemp := 15.0 + 10.0*math.Sin((dayOfYear/365.0)*2*math.Pi-math.Pi/2)
		dailyTemp := 5.0 * math.Sin((hour/24.0)*2*math.Pi-math.Pi/2)
		temperature := seasonalTemp + dailyTemp + (input-0.5)*10.0

		// Humidity inversely related to temperature
		humidity := 70.0 - temperature + (rand.Float64()-0.5)*20.0
		if humidity < 20.0 {
			humidity = 20.0
		} else if humidity > 95.0 {
			humidity = 95.0
		}

		// Pressure varies with weather systems
		pressure := 1013.25 + (input-0.5)*50.0 + (rand.Float64()-0.5)*10.0

		// Wind speed and direction
		windSpeed := math.Max(0, 10.0+input*20.0+(rand.Float64()-0.5)*5.0)
		windDirection := int(rand.Float64() * 360)

		// Weather conditions based on combined factors
		var conditions string
		if temperature < 0 {
			conditions = "snow"
		} else if humidity > 80 && temperature < 15 {
			conditions = "fog"
		} else if humidity > 70 && pressure < 1000 {
			conditions = "rain"
		} else if windSpeed > 25 {
			conditions = "windy"
		} else {
			conditions = "clear"
		}

		return WeatherData{
			StationID:     fmt.Sprintf("WX-%04d", int(input*9999)),
			Temperature:   temperature,
			Humidity:      humidity,
			Pressure:      pressure,
			WindSpeed:     windSpeed,
			WindDirection: windDirection,
			Conditions:    conditions,
			Timestamp:     timestamp.Unix(),
		}
	})

	publisher := NewConsolePublisher[WeatherData]()

	config := engine.DefaultConfig()
	config.ProductionRate = 3 * time.Second
	config.BatchSize = 1

	testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("🌤️  Starting Weather Station Example...")
	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}
}

// MarketSeeder simulates market behavior for financial metrics
type MarketSeeder struct {
	cycle float64
}

// Generate implements the Seeder interface
func (m *MarketSeeder) Generate() float64 {
	m.cycle += 0.1
	baseValue := 0.5

	// Add market cycles
	cycle := math.Sin(m.cycle*0.1) * 0.3

	// Add random market noise
	noise := (rand.Float64() - 0.5) * 0.2

	// Add trend component
	trend := math.Sin(m.cycle*0.01) * 0.2

	result := baseValue + cycle + noise + trend

	// Keep within bounds
	if result < 0 {
		result = 0
	} else if result > 1 {
		result = 1
	}

	return result
}

// Example 5: Custom Seeder with Complex Function
// Shows how to create completely custom seeder + function combinations
func CustomSeederExample() {
	type FinancialMetrics struct {
		Symbol     string  `json:"symbol"`
		Price      float64 `json:"price_usd"`
		Volume     int64   `json:"volume_24h"`
		Change     float64 `json:"change_percent_24h"`
		Volatility float64 `json:"volatility_index"`
		Trend      string  `json:"trend"`
		Timestamp  int64   `json:"timestamp_unix"`
	}

	seeder := &MarketSeeder{cycle: 0}

	// User-defined function that simulates financial metrics
	sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) FinancialMetrics {
		// Input from custom seeder represents market sentiment (0-1)
		// 0 = bear market, 1 = bull market

		// Base price varies with market sentiment
		basePrice := 100.0 + (input * 400.0) // $100-$500 range

		// Add intraday volatility
		intraday := math.Sin(float64(timestamp.Unix()%86400)*2*math.Pi/86400) * 20.0
		price := basePrice + intraday + (rand.Float64()-0.5)*10.0

		// Volume inversely related to price (higher price = lower volume)
		volume := int64((1.0-input)*1000000 + rand.Float64()*500000)

		// 24h change based on sentiment
		change := (input - 0.5) * 20.0 // -10% to +10%

		// Volatility higher during transitions
		volatility := math.Abs(input-0.5)*2.0 + rand.Float64()*0.5

		// Trend determination
		var trend string
		switch {
		case input > 0.7:
			trend = "strong_bull"
		case input > 0.6:
			trend = "bull"
		case input > 0.4:
			trend = "sideways"
		case input > 0.3:
			trend = "bear"
		default:
			trend = "strong_bear"
		}

		return FinancialMetrics{
			Symbol:     "CRYPTO-USD",
			Price:      price,
			Volume:     volume,
			Change:     change,
			Volatility: volatility,
			Trend:      trend,
			Timestamp:  timestamp.Unix(),
		}
	})

	publisher := NewConsolePublisher[FinancialMetrics]()

	config := engine.DefaultConfig()
	config.ProductionRate = 1 * time.Second
	config.BatchSize = 2

	testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("💰 Starting Custom Seeder Example...")
	if err := testEngine.Start(ctx); err != nil {
		log.Printf("Engine error: %v", err)
	}
}
