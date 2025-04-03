package dunning

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("dunning: not found")
	ErrInvalidState = errors.New("dunning: invalid state transition")
)
