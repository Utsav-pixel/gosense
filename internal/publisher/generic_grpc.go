package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Utsav-pixel/go-sensor-engine/internal/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SensorDataService defines the gRPC service interface
type SensorDataServiceClient interface {
	SendSensorData(ctx context.Context, data []byte) error
	SendSensorDataBatch(ctx context.Context, data [][]byte) error
	Close() error
}

// GenericGRPCPublisher is a generic gRPC publisher
type GenericGRPCPublisher[T any] struct {
	client SensorDataServiceClient
}

// NewGenericGRPCPublisher creates a new generic gRPC publisher
func NewGenericGRPCPublisher[T any](address string) (*GenericGRPCPublisher[T], error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := &GRPCClient{conn: conn}
	return &GenericGRPCPublisher[T]{
		client: client,
	}, nil
}

// Publish publishes a single sensor data point
func (g *GenericGRPCPublisher[T]) Publish(ctx context.Context, data engine.SensorData[T]) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return g.client.SendSensorData(ctx, payload)
}

// PublishBatch publishes a batch of sensor data points
func (g *GenericGRPCPublisher[T]) PublishBatch(ctx context.Context, data []engine.SensorData[T]) error {
	payloads := make([][]byte, len(data))
	for i, d := range data {
		payload, err := json.Marshal(d)
		if err != nil {
			return err
		}
		payloads[i] = payload
	}
	return g.client.SendSensorDataBatch(ctx, payloads)
}

// Close closes the gRPC publisher
func (g *GenericGRPCPublisher[T]) Close() error {
	return g.client.Close()
}

// GRPCClient is a simple gRPC client implementation
type GRPCClient struct {
	conn *grpc.ClientConn
}

// SendSensorData sends a single sensor data point
func (c *GRPCClient) SendSensorData(ctx context.Context, data []byte) error {
	// This is a placeholder implementation
	// In a real implementation, you would define protobuf messages and use the generated client
	fmt.Printf("Sending gRPC sensor data: %s\n", string(data))
	return nil
}

// SendSensorDataBatch sends a batch of sensor data points
func (c *GRPCClient) SendSensorDataBatch(ctx context.Context, data [][]byte) error {
	// This is a placeholder implementation
	// In a real implementation, you would define protobuf messages and use the generated client
	fmt.Printf("Sending gRPC batch of %d sensor data points\n", len(data))
	return nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
