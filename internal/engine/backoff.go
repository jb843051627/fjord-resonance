package engine

import (
	"math"
	"time"
)

type Backoff struct {
	Initial time.Duration
	Maximum time.Duration
	Factor  float64
}

func DefaultBackoff() Backoff {
	return Backoff{Initial: 100 * time.Millisecond, Maximum: 5 * time.Second, Factor: 2}
}

func (b Backoff) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if b.Initial <= 0 {
		return 0
	}
	factor := b.Factor
	if factor < 1 {
		factor = 1
	}
	delay := float64(b.Initial) * math.Pow(factor, float64(attempt))
	if time.Duration(delay) > b.Maximum && b.Maximum > 0 {
		return b.Maximum
	}
	return time.Duration(delay)
}

func (b Backoff) Schedule(start time.Time, attempts int) []time.Time {
	if attempts <= 0 {
		return nil
	}
	result := make([]time.Time, attempts)
	cursor := start
	for index := range result {
		cursor = cursor.Add(b.Delay(index))
		result[index] = cursor
	}
	return result
}
