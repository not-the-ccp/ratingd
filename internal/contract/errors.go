package contract

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("contract: not found")
	ErrInvalidState = errors.New("contract: invalid state transition")
)
