package engine

import (
	"math"
	"math/rand/v2"
	"time"
)

// TimeSeeder generates values based on time
type TimeSeeder struct {
	amplitude float64
	frequency float64
	offset    float64
}

// NewTimeSeeder creates a new time-based seeder
func NewTimeSeeder(amplitude, frequency, offset float64) *TimeSeeder {
	return &TimeSeeder{
		amplitude: amplitude,
		frequency: frequency,
		offset:    offset,
	}
}

// Generate generates a value based on current time
func (t *TimeSeeder) Generate() float64 {
	now := float64(time.Now().UnixNano()) / 1e9 // Convert to seconds with higher precision
	return t.amplitude*math.Sin(t.frequency*now) + t.offset
}

// RandomSeeder generates random values within a range
type RandomSeeder struct {
	min float64
	max float64
}

// NewRandomSeeder creates a new random seeder
func NewRandomSeeder(min, max float64) *RandomSeeder {
	return &RandomSeeder{
		min: min,
		max: max,
	}
}

// Generate generates a random value between min and max
func (r *RandomSeeder) Generate() float64 {
	return r.min + rand.Float64()*(r.max-r.min)
}

// LinearSeeder generates values that increase linearly over time
type LinearSeeder struct {
	slope  float64
	offset float64
	start  time.Time
}

// NewLinearSeeder creates a new linear seeder
func NewLinearSeeder(slope, offset float64) *LinearSeeder {
	return &LinearSeeder{
		slope:  slope,
		offset: offset,
		start:  time.Now(),
	}
}

// Generate generates a value that increases linearly
func (l *LinearSeeder) Generate() float64 {
	elapsed := float64(time.Since(l.start).Seconds())
	return l.slope*elapsed + l.offset
}

// CustomSeeder allows for custom generation functions
type CustomSeeder struct {
	generateFunc func() float64
}

// NewCustomSeeder creates a new custom seeder
func NewCustomSeeder(generateFunc func() float64) *CustomSeeder {
	return &CustomSeeder{
		generateFunc: generateFunc,
	}
}

// Generate generates a value using the custom function
func (c *CustomSeeder) Generate() float64 {
	return c.generateFunc()
}

// NormalSeeder generates values from a normal distribution
type NormalSeeder struct {
	mean   float64
	stdDev float64
}

// NewNormalSeeder creates a new normal distribution seeder
func NewNormalSeeder(mean, stdDev float64) *NormalSeeder {
	return &NormalSeeder{
		mean:   mean,
		stdDev: stdDev,
	}
}

// Generate generates a value from a normal distribution
func (n *NormalSeeder) Generate() float64 {
	return rand.NormFloat64()*n.stdDev + n.mean
}
