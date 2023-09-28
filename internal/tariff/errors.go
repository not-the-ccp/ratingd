package tariff

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("tariff: not found")
	ErrInvalidState = errors.New("tariff: invalid state transition")
)
