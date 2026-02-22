package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

// ConfigFile represents the JSON configuration file structure
type ConfigFile struct {
	Engine EngineConfig `json:"engine"`
	Seeder SeederConfig `json:"seeder"`
	Output OutputConfig `json:"output"`
}

// EngineConfig holds engine configuration
type EngineConfig struct {
	ProductionRate string `json:"production_rate"` // Duration string like "100ms", "1s"
	BatchSize      int    `json:"batch_size"`
	BatchTimeout   string `json:"batch_timeout"` // Duration string
	MaxWorkers     int    `json:"max_workers"`
}

// SeederConfig holds seeder configuration
type SeederConfig struct {
	Type     string                 `json:"type"`     // "time", "random", "linear", "normal", "custom"
	Params   map[string]interface{} `json:"params"`   // Type-specific parameters
	Function *FunctionConfig        `json:"function"` // Optional inline function definition
}

// OutputConfig holds output configuration
type OutputConfig struct {
	Type     string                 `json:"type"`     // "http", "kafka", "grpc", "console"
	Params   map[string]interface{} `json:"params"`   // Publisher-specific parameters
	Metadata map[string]string      `json:"metadata"` // Optional metadata to include in output
}

// FunctionConfig represents a simple function configuration
type FunctionConfig struct {
	Type   string                 `json:"type"`   // "simple", "lambda", "custom"
	Params map[string]interface{} `json:"params"` // Function-specific parameters
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(filename string) (*ConfigFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ConfigFile
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// ToEngineConfig converts ConfigFile to Engine Config
func (c *ConfigFile) ToEngineConfig() (Config, error) {
	productionRate, err := time.ParseDuration(c.Engine.ProductionRate)
	if err != nil {
		return Config{}, fmt.Errorf("invalid production_rate: %w", err)
	}

	batchTimeout, err := time.ParseDuration(c.Engine.BatchTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("invalid batch_timeout: %w", err)
	}

	return Config{
		ProductionRate: productionRate,
		BatchSize:      c.Engine.BatchSize,
		BatchTimeout:   batchTimeout,
		MaxWorkers:     c.Engine.MaxWorkers,
	}, nil
}

// CreateSeeder creates a seeder from configuration
func (c *ConfigFile) CreateSeeder() (Seeder, error) {
	switch c.Seeder.Type {
	case "time":
		return c.createTimeSeeder()
	case "random":
		return c.createRandomSeeder()
	case "linear":
		return c.createLinearSeeder()
	case "normal":
		return c.createNormalSeeder()
	case "custom":
		return c.createCustomSeeder()
	default:
		return nil, fmt.Errorf("unknown seeder type: %s", c.Seeder.Type)
	}
}

func (c *ConfigFile) createTimeSeeder() (Seeder, error) {
	amplitude := getFloatParam(c.Seeder.Params, "amplitude", 1.0)
	frequency := getFloatParam(c.Seeder.Params, "frequency", 0.1)
	offset := getFloatParam(c.Seeder.Params, "offset", 0.0)

	return NewTimeSeeder(amplitude, frequency, offset), nil
}

func (c *ConfigFile) createRandomSeeder() (Seeder, error) {
	min := getFloatParam(c.Seeder.Params, "min", 0.0)
	max := getFloatParam(c.Seeder.Params, "max", 1.0)

	return NewRandomSeeder(min, max), nil
}

func (c *ConfigFile) createLinearSeeder() (Seeder, error) {
	slope := getFloatParam(c.Seeder.Params, "slope", 1.0)
	offset := getFloatParam(c.Seeder.Params, "offset", 0.0)

	return NewLinearSeeder(slope, offset), nil
}

func (c *ConfigFile) createNormalSeeder() (Seeder, error) {
	mean := getFloatParam(c.Seeder.Params, "mean", 0.0)
	stdDev := getFloatParam(c.Seeder.Params, "std_dev", 1.0)

	return NewNormalSeeder(mean, stdDev), nil
}

func (c *ConfigFile) createCustomSeeder() (Seeder, error) {
	// For custom seeders, we'd need to load Go code or use a scripting language
	// For now, return a simple sine wave as example
	return NewCustomSeeder(func() float64 {
		t := float64(time.Now().UnixNano()) / 1e9
		return getFloatParam(c.Seeder.Params, "amplitude", 1.0) *
			(0.3*math.Sin(t*2.0) + 0.2*math.Sin(t*7.3) + 0.1*math.Sin(t*13.7))
	}), nil
}

// Helper functions for parameter extraction
func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if parsed, err := parseFloat(v); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}

func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if parsed, err := parseInt(v); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}

func getStringParam(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// Parse functions for string parameters
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// CreateEngineFromConfig creates a complete engine configuration from file
func CreateEngineFromConfig[T any](filename string, function SensorFunction[T], publisher Publisher[T]) (*Engine[T], error) {
	configFile, err := LoadConfigFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	engineConfig, err := configFile.ToEngineConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert engine config: %w", err)
	}

	seeder, err := configFile.CreateSeeder()
	if err != nil {
		return nil, fmt.Errorf("failed to create seeder: %w", err)
	}

	return NewEngine(engineConfig, seeder, function, publisher), nil
}

// DefaultConfigFile returns a default configuration structure
func DefaultConfigFile() *ConfigFile {
	return &ConfigFile{
		Engine: EngineConfig{
			ProductionRate: "100ms",
			BatchSize:      100,
			BatchTimeout:   "500ms",
			MaxWorkers:     3,
		},
		Seeder: SeederConfig{
			Type: "time",
			Params: map[string]interface{}{
				"amplitude": 1.0,
				"frequency": 0.1,
				"offset":    0.0,
			},
		},
		Output: OutputConfig{
			Type:   "console",
			Params: map[string]interface{}{},
			Metadata: map[string]string{
				"version": "1.0",
				"source":  "go-sensor-engine",
			},
		},
	}
}

// SaveConfigToFile saves configuration to a JSON file
func SaveConfigToFile(config *ConfigFile, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
