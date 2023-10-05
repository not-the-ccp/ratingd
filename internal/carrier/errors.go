package carrier

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("carrier: not found")
	ErrInvalidState     = errors.New("carrier: invalid state transition")
	ErrCarrierEmptyCode = errors.New("carrier: carrier rejected, empty code")
)
