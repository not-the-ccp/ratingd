package invoice

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("invoice: not found")
	ErrInvalidState = errors.New("invoice: invalid state transition")
)
