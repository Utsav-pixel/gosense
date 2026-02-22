# 🎯 **SEEDER + FUNCTION INTEGRATION GUIDE**

## 📚 **Understanding the Relationship**

In the Generic Sensor Engine, **seeders** and **functions** work together to create realistic sensor data:

- **Seeder**: Generates input values (0-1) that drive your sensor logic
- **Function**: Takes seeder input + timestamp → produces your sensor data

This separation allows you to:
1. Control **how** data changes over time (seeder)
2. Control **what** data gets generated (function)

---

## 🌱 **Available Seeders**

### 1. **TimeSeeder** - Time-based patterns
```go
// Generates values that change over time with sine waves
seeder := engine.NewTimeSeeder(amplitude, frequency, offset)

// Example: Daily temperature cycles
seeder := engine.NewTimeSeeder(1.0, 0.1, 20.0)
// amplitude=1.0: temperature varies ±1°C
// frequency=0.1: slow daily cycle  
// offset=20.0: base temperature 20°C
```

### 2. **RandomSeeder** - Random events
```go
// Generates random values in a range
seeder := engine.NewRandomSeeder(min, max)

// Example: Random IoT device states
seeder := engine.NewRandomSeeder(0.0, 1.0)
// Generates random values between 0 and 1
```

### 3. **LinearSeeder** - Progressive changes
```go
// Generates values that increase/decrease linearly
seeder := engine.NewLinearSeeder(increment, start)

// Example: Machine wear over time
seeder := engine.NewLinearSeeder(0.01, 0.1)
// increment=0.01: wear increases by 0.01 each generation
// start=0.1: starts at 10% wear
```

### 4. **NormalSeeder** - Natural distributions
```go
// Generates values following normal distribution
seeder := engine.NewNormalSeeder(mean, stdDev)

// Example: Natural weather variations
seeder := engine.NewNormalSeeder(0.5, 0.2)
// mean=0.5: centered around middle
// stdDev=0.2: most values within 0.1-0.9 range
```

### 5. **Custom Seeder** - Your own logic
```go
// Create your own seeder by implementing the Seeder interface
type MarketSeeder struct {
    cycle float64
}

func (m *MarketSeeder) Generate() float64 {
    m.cycle += 0.1
    // Your custom logic here
    return math.Sin(m.cycle) * 0.5 + 0.5
}

seeder := &MarketSeeder{cycle: 0}
```

---

## 🔧 **Function Types**

### 1. **BasicSensorFunction** - Simple transformations
```go
sensorFunc := engine.NewBasicSensorFunction(func(input float64, timestamp time.Time) YourDataType {
    // Your logic here
    return yourData
})
```

### 2. **Function** - User-defined logic
```go
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) YourDataType {
    // Your logic here
    return yourData
})
```

### 3. **LambdaSensorFunction** - Inline anonymous functions
```go
sensorFunc := engine.NewLambdaSensorFunction(func(input float64, timestamp time.Time) YourDataType {
    // Your logic here
    return yourData
})
```

---

## 🎯 **Real-World Integration Examples**

### **Example 1: Temperature Sensor with Time-based Seeder**

```go
type TemperatureReading struct {
    Celsius    float64 `json:"celsius"`
    Fahrenheit float64 `json:"fahrenheit"`
    Humidity   float64 `json:"humidity_percent"`
    Location   string  `json:"location"`
}

// Time-based seeder simulates daily temperature cycles
seeder := engine.NewTimeSeeder(1.0, 0.1, 20.0)

// Function uses seeder input + time for realistic temperature data
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) TemperatureReading {
    // Input represents environmental factor (0-1) from seeder
    baseTemp := input
    
    // Add diurnal pattern (daily cycle)
    hour := float64(timestamp.Hour()) + float64(timestamp.Minute())/60.0
    radian := (hour / 24.0) * 2 * math.Pi
    diurnal := 5.0 * math.Sin(radian - math.Pi/2)
    
    // Add random noise
    noise := (rand.Float64() - 0.5) * 1.0
    
    celsius := baseTemp + diurnal + noise
    fahrenheit := celsius*9/5 + 32
    
    // Humidity inversely related to temperature
    humidity := 70.0 - celsius
    if humidity < 30.0 { humidity = 30.0 } else if humidity > 90.0 { humidity = 90.0 }
    
    return TemperatureReading{
        Celsius:    celsius,
        Fahrenheit: fahrenheit,
        Humidity:   humidity,
        Location:   "Server Room A",
    }
})
```

**How it works:**
- **Seeder**: Generates values that cycle over time (simulating day/night)
- **Function**: Uses seeder input as base temperature + adds daily patterns + noise

### **Example 2: IoT Device with Random Seeder**

```go
type IoTReading struct {
    DeviceID    string  `json:"device_id"`
    Battery     float64 `json:"battery_percent"`
    Signal      int     `json:"signal_strength_dbm"`
    Temperature float64 `json:"temperature_celsius"`
    Status      string  `json:"status"`
}

// Random seeder simulates random device states
seeder := engine.NewRandomSeeder(0.0, 1.0)

// Function maps random input to device metrics
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) IoTReading {
    // Input represents device stress/activity level (0-1)
    
    // Battery decreases with activity
    battery := 100.0 - (input * 30.0) // 70-100% range
    
    // Signal varies with activity
    signal := -30 - int(input*40) // -30 to -70 dBm
    
    // Temperature increases with activity
    temperature := 25.0 + (input * 15.0) // 25-40°C
    
    // Status based on battery
    var status string
    switch {
    case battery > 80: status = "excellent"
    case battery > 50: status = "good"
    case battery > 20: status = "low"
    default: status = "critical"
    }
    
    return IoTReading{
        DeviceID:    fmt.Sprintf("iot-%04d", int(input*9999)),
        Battery:     battery,
        Signal:      signal,
        Temperature: temperature,
        Status:      status,
    }
})
```

**How it works:**
- **Seeder**: Generates random device activity levels
- **Function**: Maps activity to realistic device metrics (battery drain, signal loss, etc.)

### **Example 3: Industrial Machine with Linear Seeder**

```go
type MachineMetrics struct {
    MachineID  string  `json:"machine_id"`
    Vibration  float64 `json:"vibration_mm_s"`
    Pressure   float64 `json:"pressure_bar"`
    RPM        int     `json:"rpm"`
    Temperature float64 `json:"temperature_celsius"`
    Efficiency  float64 `json:"efficiency_percent"`
    Status     string  `json:"status"`
}

// Linear seeder simulates gradual machine wear
seeder := engine.NewLinearSeeder(0.01, 0.1)

// Function simulates machine degradation over time
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) MachineMetrics {
    // Input represents machine wear factor (increases over time)
    
    // Vibration increases with wear
    vibration := input * 8.0 // 0.8 to 8+ mm/s
    
    // Pressure becomes unstable with wear
    pressure := 2.5 + (input * 1.5) + (rand.Float64()-0.5)*0.5
    
    // RPM decreases as machine wears out
    rpm := 1800 - int(input * 400) // Up to 400 RPM reduction
    
    // Temperature increases with friction
    temperature := 25.0 + (input * 30.0)
    
    // Efficiency decreases with wear
    efficiency := 100.0 - (input * 40.0) // 60% to 100%
    
    // Status based on wear level
    var status string
    switch {
    case input > 0.8: status = "critical_maintenance"
    case input > 0.6: status = "warning"
    case input > 0.3: status = "monitor"
    default: status = "normal"
    }
    
    return MachineMetrics{
        MachineID:  fmt.Sprintf("CNC-%03d", int(input*999)),
        Vibration:  vibration,
        Pressure:   pressure,
        RPM:        rpm,
        Temperature: temperature,
        Efficiency:  efficiency,
        Status:     status,
    }
})
```

**How it works:**
- **Seeder**: Linearly increasing wear factor (0.1 → 0.11 → 0.12...)
- **Function**: Maps wear to realistic machine degradation patterns

### **Example 4: Weather Station with Normal Seeder**

```go
type WeatherData struct {
    StationID    string  `json:"station_id"`
    Temperature  float64 `json:"temperature_celsius"`
    Humidity     float64 `json:"humidity_percent"`
    Pressure     float64 `json:"pressure_hpa"`
    WindSpeed    float64 `json:"wind_speed_kmh"`
    Conditions   string  `json:"conditions"`
}

// Normal seeder simulates natural weather variations
seeder := engine.NewNormalSeeder(0.5, 0.2)

// Function creates realistic weather patterns
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) WeatherData {
    // Input represents weather variability (normally distributed)
    
    // Temperature with seasonal + daily variations
    hour := float64(timestamp.Hour())
    dayOfYear := float64(timestamp.YearDay())
    
    seasonalTemp := 15.0 + 10.0*math.Sin((dayOfYear/365.0)*2*math.Pi-math.Pi/2)
    dailyTemp := 5.0*math.Sin((hour/24.0)*2*math.Pi-math.Pi/2)
    temperature := seasonalTemp + dailyTemp + (input-0.5)*10.0
    
    // Humidity inversely related to temperature
    humidity := 70.0 - temperature + (rand.Float64()-0.5)*20.0
    if humidity < 20.0 { humidity = 20.0 } else if humidity > 95.0 { humidity = 95.0 }
    
    // Pressure varies with weather systems
    pressure := 1013.25 + (input-0.5)*50.0 + (rand.Float64()-0.5)*10.0
    
    // Wind speed
    windSpeed := math.Max(0, 10.0+input*20.0+(rand.Float64()-0.5)*5.0)
    
    // Weather conditions
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
        StationID:    fmt.Sprintf("WX-%04d", int(input*9999)),
        Temperature:  temperature,
        Humidity:     humidity,
        Pressure:     pressure,
        WindSpeed:    windSpeed,
        Conditions:   conditions,
    }
})
```

**How it works:**
- **Seeder**: Normally distributed weather variability (most values near 0.5)
- **Function**: Creates realistic weather with seasonal patterns, daily cycles, and conditions

---

## 🚀 **Best Practices**

### **1. Choose the Right Seeder**
- **Time patterns**: Use `TimeSeeder` for cyclical data (temperature, daily cycles)
- **Random events**: Use `RandomSeeder` for unpredictable events (IoT devices, failures)
- **Progressive changes**: Use `LinearSeeder` for degradation/growth (machine wear, learning)
- **Natural phenomena**: Use `NormalSeeder` for realistic distributions (weather, human behavior)
- **Complex patterns**: Use `Custom Seeder` for specific business logic

### **2. Design Your Function Logic**
```go
// Good: Clear separation of concerns
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) YourData {
    // 1. Use seeder input as primary driver
    baseValue := input * scaleFactor
    
    // 2. Add time-based patterns
    timeEffect := calculateTimeEffect(timestamp)
    
    // 3. Add realistic noise
    noise := (rand.Float64() - 0.5) * noiseLevel
    
    // 4. Apply business logic
    result := applyBusinessLogic(baseValue + timeEffect + noise)
    
    return result
})
```

### **3. Handle Edge Cases**
```go
sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) YourData {
    // Clamp values to realistic ranges
    if result < minValue { result = minValue }
    if result > maxValue { result = maxValue }
    
    // Handle special conditions
    if timestamp.Hour() < 6 {
        // Night-time behavior
    }
    
    return result
})
```

---

## 🎯 **Quick Start Template**

```go
package main

import (
    "time"
    "github.com/Utsav-pixel/go-sensor-engine/internal/engine"
)

func main() {
    // 1. Define your data structure
    type MySensorData struct {
        Value    float64 `json:"value"`
        Status   string  `json:"status"`
        Location string  `json:"location"`
    }
    
    // 2. Choose your seeder
    seeder := engine.NewTimeSeeder(1.0, 0.1, 0.0)
    
    // 3. Create your function
    sensorFunc := engine.NewFunction(func(input float64, timestamp time.Time) MySensorData {
        value := input * 100.0
        
        var status string
        if value > 80 {
            status = "high"
        } else if value > 50 {
            status = "medium"
        } else {
            status = "low"
        }
        
        return MySensorData{
            Value:    value,
            Status:   status,
            Location: "sensor-001",
        }
    })
    
    // 4. Create publisher and engine
    publisher := NewYourPublisher[MySensorData]()
    config := engine.DefaultConfig()
    
    testEngine := engine.NewEngine(config, seeder, sensorFunc, publisher)
    
    // 5. Run it
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    testEngine.Start(ctx)
}
```

---

## 🏆 **Key Takeaways**

1. **Seeders control the "how"** - how your data changes over time
2. **Functions control the "what"** - what data gets generated
3. **Separation of concerns** makes code more maintainable and testable
4. **Choose the right seeder** for your domain (time-based, random, linear, normal, custom)
5. **Add realistic patterns** - noise, cycles, business logic
6. **Handle edge cases** - clamp values, special conditions

This architecture gives you **complete flexibility** while maintaining **clean, testable code**! 🚀
