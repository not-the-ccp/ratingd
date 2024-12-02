package settlement

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("settlement: not found")
	ErrInvalidState = errors.New("settlement: invalid state transition")
)
