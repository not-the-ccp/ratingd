package clock

import (
	"time"
)

// Clock abstracts the current time so services can be tested.
type Clock interface {
	Now() time.Time
}

// System is the wall clock.
type System struct {
}

func (s System) Now() time.Time {
	return time.Now().UTC()
}
