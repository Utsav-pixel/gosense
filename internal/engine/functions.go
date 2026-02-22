package engine

import (
	"time"
)

// BasicSensorFunction provides a simple implementation for basic sensor data generation
type BasicSensorFunction[T any] struct {
	transformFunc func(float64, time.Time) T
}

// NewBasicSensorFunction creates a new basic sensor function with a custom transform function
func NewBasicSensorFunction[T any](transformFunc func(float64, time.Time) T) *BasicSensorFunction[T] {
	return &BasicSensorFunction[T]{
		transformFunc: transformFunc,
	}
}

// Generate creates sensor data using the transform function
func (f *BasicSensorFunction[T]) Generate(input float64, timestamp time.Time) T {
	return f.transformFunc(input, timestamp)
}

// Function allows users to define their own sensor data generation logic
type Function[T any] struct {
	generateFunc func(float64, time.Time) T
}

// NewFunction creates a new user-defined sensor function
func NewFunction[T any](generateFunc func(float64, time.Time) T) *Function[T] {
	return &Function[T]{
		generateFunc: generateFunc,
	}
}

// Generate creates sensor data using the user-defined function
func (f *Function[T]) Generate(input float64, timestamp time.Time) T {
	return f.generateFunc(input, timestamp)
}

// LambdaSensorFunction provides a simple function wrapper for inline usage
type LambdaSensorFunction[T any] struct {
	lambda func(float64, time.Time) T
}

// NewLambdaSensorFunction creates a sensor function from a lambda/anonymous function
func NewLambdaSensorFunction[T any](lambda func(float64, time.Time) T) *LambdaSensorFunction[T] {
	return &LambdaSensorFunction[T]{
		lambda: lambda,
	}
}

// Generate generates sensor data using the lambda function
func (l *LambdaSensorFunction[T]) Generate(input float64, timestamp time.Time) T {
	return l.lambda(input, timestamp)
}
