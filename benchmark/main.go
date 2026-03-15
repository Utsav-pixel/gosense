package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/Utsav-pixel/gosense"
)

// Benchmark data structure
type BenchmarkData struct {
	ID        string  `json:"id"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Location  string  `json:"location"`
	Status    string  `json:"status"`
}

func main() {
	log.Println("=== GoSense Performance Benchmark ===")

	// Get machine specs
	printMachineSpecs()

	// Test 1: Data Generation Speed
	log.Println("\n📊 Test 1: Data Generation Speed")
	dataCount, elapsed := benchmarkDataGeneration()
	throughput := float64(dataCount) / elapsed.Seconds()
	log.Printf("  Generated: %d data points", dataCount)
	log.Printf("  Time: %.2f seconds", elapsed.Seconds())
	log.Printf("  Throughput: %.0f data points/second", throughput)

	// Test 2: Memory Usage
	log.Println("\n💾 Test 2: Memory Usage")
	memoryUsed, dataCount2, elapsed2 := benchmarkMemoryUsage()
	memoryPer100K := float64(memoryUsed) * 100000 / float64(dataCount2)
	log.Printf("  Data points: %d", dataCount2)
	log.Printf("  Time: %.2f seconds", elapsed2.Seconds())
	log.Printf("  Memory used: %d MB", memoryUsed)
	log.Printf("  Memory per 100K points: %.2f MB", memoryPer100K)

	// Test 3: Latency
	log.Println("\n⚡ Test 3: Latency Measurement")
	avgLatencyMicros := benchmarkLatency()
	log.Printf("  Average latency: %.1f microseconds per data point", avgLatencyMicros)
	log.Printf("  Average latency: %.3f ms per data point", avgLatencyMicros/1000)

	// Summary
	log.Println("\n🎯 BENCHMARK SUMMARY")
	log.Printf("  Machine: %s/%s, %d cores, Go %s", runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Version())
	log.Printf("  Throughput: %.0f data points/second", throughput)
	log.Printf("  Memory: %.2f MB per 100K points", memoryPer100K)
	log.Printf("  Latency: %.1f microseconds per point", avgLatencyMicros)
}

func printMachineSpecs() {
	log.Println("🖥️  Machine Specifications:")
	log.Printf("  OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Printf("  Go Version: %s", runtime.Version())
	log.Printf("  CPU Cores: %d", runtime.NumCPU())
	log.Printf("  GOMAXPROCS: %d", runtime.GOMAXPROCS(0))
}

func benchmarkDataGeneration() (int, time.Duration) {
	config := gosense.HighThroughputConfig()
	config.ProductionRate = 1 * time.Millisecond // Fast as possible
	config.BatchSize = 10000
	config.MaxWorkers = runtime.NumCPU()

	seeder := gosense.NewRandomSeeder(0.0, 1.0)
	sensorFunc := gosense.NewFunction(func(input float64, timestamp time.Time) BenchmarkData {
		return BenchmarkData{
			ID:        fmt.Sprintf("sensor-%d", int(input*1000000)),
			Value:     input * 1000,
			Timestamp: timestamp.UnixNano(),
			Location:  "benchmark-zone",
			Status:    "active",
		}
	})

	publisher := &SilentPublisher{}
	engine := gosense.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	if err := engine.Start(ctx); err != nil {
		log.Printf("  Error: %v", err)
		return 0, 0
	}
	elapsed := time.Since(start)

	return publisher.count, elapsed
}

func benchmarkMemoryUsage() (int, int, time.Duration) {
	config := gosense.HighThroughputConfig()
	config.ProductionRate = 1 * time.Millisecond
	config.BatchSize = 5000
	config.MaxWorkers = runtime.NumCPU()

	seeder := gosense.NewRandomSeeder(0.0, 1.0)
	sensorFunc := gosense.NewFunction(func(input float64, timestamp time.Time) BenchmarkData {
		return BenchmarkData{
			ID:        fmt.Sprintf("mem-test-%d", int(input*100000)),
			Value:     input * 500,
			Timestamp: timestamp.UnixNano(),
			Location:  "memory-test",
			Status:    "testing",
		}
	})

	publisher := &SilentPublisher{}
	engine := gosense.NewEngine(config, seeder, sensorFunc, publisher)

	// Measure memory before
	var beforeMem runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&beforeMem)

	// Run for 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	if err := engine.Start(ctx); err != nil {
		log.Printf("  Error: %v", err)
		return 0, 0, 0
	}
	elapsed := time.Since(start)

	// Measure memory after
	var afterMem runtime.MemStats
	runtime.ReadMemStats(&afterMem)

	memoryUsed := int((afterMem.Alloc - beforeMem.Alloc) / 1024 / 1024) // MB

	return memoryUsed, publisher.count, elapsed
}

func benchmarkLatency() float64 {
	config := gosense.LowLatencyConfig()
	config.ProductionRate = 1 * time.Millisecond
	config.BatchSize = 1
	config.MaxWorkers = 1

	seeder := gosense.NewRandomSeeder(0.0, 1.0)
	sensorFunc := gosense.NewFunction(func(input float64, timestamp time.Time) BenchmarkData {
		return BenchmarkData{
			ID:        fmt.Sprintf("latency-%d", int(input*1000)),
			Value:     input * 10,
			Timestamp: timestamp.UnixNano(),
			Location:  "latency-test",
			Status:    "measuring",
		}
	})

	publisher := &SilentPublisher{}
	engine := gosense.NewEngine(config, seeder, sensorFunc, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	if err := engine.Start(ctx); err != nil {
		log.Printf("  Error: %v", err)
		return 0
	}
	elapsed := time.Since(start)

	if publisher.count > 0 {
		avgLatencyMicros := float64(elapsed.Nanoseconds()) / float64(publisher.count) / 1000
		return avgLatencyMicros
	}
	return 0
}

// SilentPublisher doesn't output anything, just counts
type SilentPublisher struct {
	count int
}

func (p *SilentPublisher) Publish(ctx context.Context, data gosense.SensorData[BenchmarkData]) error {
	p.count++
	return nil
}

func (p *SilentPublisher) PublishBatch(ctx context.Context, data []gosense.SensorData[BenchmarkData]) error {
	p.count += len(data)
	return nil
}

func (p *SilentPublisher) Close() error {
	return nil
}
