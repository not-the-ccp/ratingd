package export

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("export: not found")
	ErrInvalidState = errors.New("export: invalid state transition")
)
