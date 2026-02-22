package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
)

// GenericHTTPPublisher is a generic HTTP publisher
type GenericHTTPPublisher[T any] struct {
	endpoint string
	client   *http.Client
}

// NewGenericHTTPPublisher creates a new generic HTTP publisher
func NewGenericHTTPPublisher[T any](endpoint string) *GenericHTTPPublisher[T] {
	return &GenericHTTPPublisher[T]{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Publish publishes a single sensor data point
func (h *GenericHTTPPublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// PublishBatch publishes a batch of sensor data points
func (h *GenericHTTPPublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Close closes the HTTP publisher
func (h *GenericHTTPPublisher[T]) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}
