package engine

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestConfigFile_LoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	configData := `{
		"engine": {
			"production_rate": "100ms",
			"batch_size": 50,
			"batch_timeout": "500ms",
			"max_workers": 3
		},
		"seeder": {
			"type": "time",
			"params": {
				"amplitude": 1.0,
				"frequency": 0.1,
				"offset": 0.0
			}
		},
		"output": {
			"type": "http",
			"params": {
				"endpoint": "https://api.example.com/data"
			},
			"metadata": {
				"version": "1.0"
			}
		}
	}`

	tmpFile, err := os.CreateTemp("", "test-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configData); err != nil {
		t.Fatalf("Failed to write config data: %v", err)
	}
	tmpFile.Close()

	// Test loading config
	config, err := LoadConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test engine config
	if config.Engine.ProductionRate != "100ms" {
		t.Errorf("Expected production_rate '100ms', got '%s'", config.Engine.ProductionRate)
	}
	if config.Engine.BatchSize != 50 {
		t.Errorf("Expected batch_size 50, got %d", config.Engine.BatchSize)
	}
	if config.Engine.MaxWorkers != 3 {
		t.Errorf("Expected max_workers 3, got %d", config.Engine.MaxWorkers)
	}

	// Test seeder config
	if config.Seeder.Type != "time" {
		t.Errorf("Expected seeder type 'time', got '%s'", config.Seeder.Type)
	}
	if getFloatParam(config.Seeder.Params, "amplitude", 0.0) != 1.0 {
		t.Errorf("Expected amplitude 1.0, got %f", getFloatParam(config.Seeder.Params, "amplitude", 0.0))
	}

	// Test output config
	if config.Output.Type != "http" {
		t.Errorf("Expected output type 'http', got '%s'", config.Output.Type)
	}
	if getStringParam(config.Output.Params, "endpoint", "") != "https://api.example.com/data" {
		t.Errorf("Expected endpoint 'https://api.example.com/data', got '%s'", getStringParam(config.Output.Params, "endpoint", ""))
	}
}

func TestConfigFile_ToEngineConfig(t *testing.T) {
	config := &ConfigFile{
		Engine: EngineConfig{
			ProductionRate: "100ms",
			BatchSize:      50,
			BatchTimeout:   "500ms",
			MaxWorkers:     3,
		},
	}

	engineConfig, err := config.ToEngineConfig()
	if err != nil {
		t.Fatalf("Failed to convert engine config: %v", err)
	}

	expectedProductionRate := 100 * time.Millisecond
	if engineConfig.ProductionRate != expectedProductionRate {
		t.Errorf("Expected production rate %v, got %v", expectedProductionRate, engineConfig.ProductionRate)
	}

	if engineConfig.BatchSize != 50 {
		t.Errorf("Expected batch size 50, got %d", engineConfig.BatchSize)
	}

	expectedBatchTimeout := 500 * time.Millisecond
	if engineConfig.BatchTimeout != expectedBatchTimeout {
		t.Errorf("Expected batch timeout %v, got %v", expectedBatchTimeout, engineConfig.BatchTimeout)
	}

	if engineConfig.MaxWorkers != 3 {
		t.Errorf("Expected max workers 3, got %d", engineConfig.MaxWorkers)
	}
}

func TestConfigFile_CreateSeeder(t *testing.T) {
	tests := []struct {
		name        string
		seederType  string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name:       "TimeSeeder",
			seederType: "time",
			params: map[string]interface{}{
				"amplitude": 1.0,
				"frequency": 0.1,
				"offset":    0.0,
			},
			expectError: false,
		},
		{
			name:       "RandomSeeder",
			seederType: "random",
			params: map[string]interface{}{
				"min": 0.0,
				"max": 1.0,
			},
			expectError: false,
		},
		{
			name:       "LinearSeeder",
			seederType: "linear",
			params: map[string]interface{}{
				"slope":  1.0,
				"offset": 0.0,
			},
			expectError: false,
		},
		{
			name:       "NormalSeeder",
			seederType: "normal",
			params: map[string]interface{}{
				"mean":    0.0,
				"std_dev": 1.0,
			},
			expectError: false,
		},
		{
			name:       "CustomSeeder",
			seederType: "custom",
			params: map[string]interface{}{
				"amplitude": 2.0,
			},
			expectError: false,
		},
		{
			name:        "InvalidSeeder",
			seederType:  "invalid",
			params:      map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ConfigFile{
				Seeder: SeederConfig{
					Type:   tt.seederType,
					Params: tt.params,
				},
			}

			seeder, err := config.CreateSeeder()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for seeder type '%s'", tt.seederType)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error creating seeder: %v", err)
				return
			}

			if seeder == nil {
				t.Error("Expected non-nil seeder")
				return
			}

			// Test that seeder generates values
			value := seeder.Generate()
			if !isFinite(value) {
				t.Errorf("Expected finite value, got %v", value)
			}
		})
	}
}

func TestDefaultConfigFile(t *testing.T) {
	config := DefaultConfigFile()

	if config.Engine.ProductionRate != "100ms" {
		t.Errorf("Expected default production_rate '100ms', got '%s'", config.Engine.ProductionRate)
	}
	if config.Engine.BatchSize != 100 {
		t.Errorf("Expected default batch_size 100, got %d", config.Engine.BatchSize)
	}
	if config.Engine.MaxWorkers != 3 {
		t.Errorf("Expected default max_workers 3, got %d", config.Engine.MaxWorkers)
	}

	if config.Seeder.Type != "time" {
		t.Errorf("Expected default seeder type 'time', got '%s'", config.Seeder.Type)
	}

	if config.Output.Type != "console" {
		t.Errorf("Expected default output type 'console', got '%s'", config.Output.Type)
	}
}

func TestSaveConfigToFile(t *testing.T) {
	config := DefaultConfigFile()

	tmpFile, err := os.CreateTemp("", "test-save-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Save config
	err = SaveConfigToFile(config, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load and verify
	loadedConfig, err := LoadConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedConfig.Engine.ProductionRate != config.Engine.ProductionRate {
		t.Errorf("Saved/loaded production_rate mismatch: expected '%s', got '%s'",
			config.Engine.ProductionRate, loadedConfig.Engine.ProductionRate)
	}

	if loadedConfig.Seeder.Type != config.Seeder.Type {
		t.Errorf("Saved/loaded seeder type mismatch: expected '%s', got '%s'",
			config.Seeder.Type, loadedConfig.Seeder.Type)
	}
}

func TestCreateEngineFromConfig(t *testing.T) {
	// Create a temporary config file
	configData := `{
		"engine": {
			"production_rate": "50ms",
			"batch_size": 10,
			"batch_timeout": "100ms",
			"max_workers": 2
		},
		"seeder": {
			"type": "random",
			"params": {
				"min": 0.0,
				"max": 1.0
			}
		},
		"output": {
			"type": "http",
			"params": {
				"endpoint": "https://example.com"
			}
		}
	}`

	tmpFile, err := os.CreateTemp("", "test-engine-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configData); err != nil {
		t.Fatalf("Failed to write config data: %v", err)
	}
	tmpFile.Close()

	// Mock publisher for testing
	mockPublisher := &mockTestPublisher[float64]{}

	// Create engine from config
	engine, err := CreateEngineFromConfig(tmpFile.Name(),
		NewLambdaSensorFunction(func(input float64, timestamp time.Time) float64 {
			return input * 2.0
		}), mockPublisher)

	if err != nil {
		t.Fatalf("Failed to create engine from config: %v", err)
	}

	if engine == nil {
		t.Fatal("Expected non-nil engine")
	}
}

// Helper functions and mocks
func isFinite(f float64) bool {
	return !(f != f || f > 1.797693134862315708145274237317043567981e+308 || f < -1.797693134862315708145274237317043567981e+308)
}

type mockTestPublisher[T any] struct{}

func (m *mockTestPublisher[T]) Publish(ctx context.Context, data SensorData[T]) error {
	return nil
}

func (m *mockTestPublisher[T]) PublishBatch(ctx context.Context, data []SensorData[T]) error {
	return nil
}

func (m *mockTestPublisher[T]) Close() error {
	return nil
}
