package engine

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// Start starts the sensor engine and returns an error if any
func (e *Engine[T]) Start(ctx context.Context) error {
	// Create channels for data flow
	dataChan := make(chan SensorData[T], 100)
	batchChan := make(chan []SensorData[T], 10)

	// Wait groups for graceful shutdown
	var dataWG, batchWG, publishWG sync.WaitGroup

	// Start data generator
	dataWG.Add(1)
	go e.generateData(ctx, dataChan, &dataWG)

	// Start batch processor
	batchWG.Add(1)
	go e.processBatches(ctx, dataChan, batchChan, &batchWG)

	// Start publisher workers
	for i := 0; i < e.config.MaxWorkers; i++ {
		publishWG.Add(1)
		go e.publishWorker(ctx, batchChan, &publishWG)
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Wait for data generator to finish first
	dataWG.Wait()

	// Then close data channel to signal batch processor to stop
	close(dataChan)

	// Wait for batch processor to finish
	batchWG.Wait()

	// Close batch channel to signal publisher workers to stop
	close(batchChan)

	// Wait for publisher workers to finish
	publishWG.Wait()

	// Close publisher
	if err := e.publisher.Close(); err != nil {
		return fmt.Errorf("error closing publisher: %w", err)
	}

	return nil
}

// generateData continuously generates sensor data
func (e *Engine[T]) generateData(ctx context.Context, dataChan chan<- SensorData[T], wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(e.config.ProductionRate)
	defer ticker.Stop()

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			input := e.seeder.Generate()
			timestamp := time.Now()
			data := e.function.Generate(input, timestamp)

			sensorData := SensorData[T]{
				ID:        fmt.Sprintf("sensor-%d", counter),
				Timestamp: timestamp,
				Data:      data,
				Quality:   determineQuality(),
			}

			select {
			case dataChan <- sensorData:
				counter++
			case <-ctx.Done():
				return
			}
		}
	}
}

// processBatches collects data into batches and sends them to batch channel
func (e *Engine[T]) processBatches(ctx context.Context, dataChan <-chan SensorData[T], batchChan chan<- []SensorData[T], wg *sync.WaitGroup) {
	defer wg.Done()

	batch := make([]SensorData[T], 0, e.config.BatchSize)
	batchTicker := time.NewTicker(e.config.BatchTimeout)
	defer batchTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Send remaining batch before exiting
			if len(batch) > 0 {
				select {
				case batchChan <- batch:
				case <-ctx.Done():
				}
			}
			return

		case data, ok := <-dataChan:
			if !ok {
				// Data channel closed, send remaining batch and exit
				if len(batch) > 0 {
					select {
					case batchChan <- batch:
					case <-ctx.Done():
					}
				}
				return
			}

			batch = append(batch, data)

			// Send batch if it reaches the size limit
			if len(batch) >= e.config.BatchSize {
				select {
				case batchChan <- batch:
					batch = make([]SensorData[T], 0, e.config.BatchSize)
				case <-ctx.Done():
					return
				}
			}

		case <-batchTicker.C:
			// Send batch if it has data and timeout is reached
			if len(batch) > 0 {
				select {
				case batchChan <- batch:
					batch = make([]SensorData[T], 0, e.config.BatchSize)
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// publishWorker publishes batches to the configured publisher
func (e *Engine[T]) publishWorker(ctx context.Context, batchChan <-chan []SensorData[T], wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-batchChan:
			if !ok {
				return
			}

			if err := e.publisher.PublishBatch(ctx, batch); err != nil {
				// Log error but continue processing
				fmt.Printf("Error publishing batch: %v\n", err)
			}
		}
	}
}

// determineQuality randomly determines the quality of sensor data
func determineQuality() Quality {
	r := rand.Float64()
	switch {
	case r < 0.01:
		return QualityCorrupt
	case r < 0.03:
		return QualityPartial
	case r < 0.08:
		return QualityNoisy
	default:
		return QualityOK
	}
}

// DefaultConfig returns a default engine configuration
func DefaultConfig() Config {
	return Config{
		ProductionRate: 100 * time.Millisecond,
		BatchSize:      100,
		BatchTimeout:   500 * time.Millisecond,
		MaxWorkers:     3,
	}
}

// HighThroughputConfig returns a configuration optimized for high throughput
func HighThroughputConfig() Config {
	return Config{
		ProductionRate: 10 * time.Millisecond,
		BatchSize:      1000,
		BatchTimeout:   100 * time.Millisecond,
		MaxWorkers:     10,
	}
}

// LowLatencyConfig returns a configuration optimized for low latency
func LowLatencyConfig() Config {
	return Config{
		ProductionRate: 50 * time.Millisecond,
		BatchSize:      10,
		BatchTimeout:   25 * time.Millisecond,
		MaxWorkers:     5,
	}
}
