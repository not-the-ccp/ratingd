package rating

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("rating: not found")
	ErrInvalidState = errors.New("rating: invalid state transition")
)
