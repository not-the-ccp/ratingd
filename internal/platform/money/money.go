package money

import (
	"fmt"
	"math"
)

// Amount is a monetary value in minor units (cents).
type Amount int64

// Add returns the sum of a and other.
func (a Amount) Add(other Amount) Amount {
	return a + other
}

func (a Amount) Sub(other Amount) Amount {
	return a - other
}

// Percent returns the given percentage of a, rounded half away from zero.
func (a Amount) Percent(percent float64) Amount {
	return Amount(math.Round(float64(a) * percent / 100))
}

func (a Amount) IsZero() bool {
	return a == 0
}

func (a Amount) String() string {
	return fmt.Sprintf("%d.%02d", int64(a)/100, abs(int64(a))%100)
}

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
