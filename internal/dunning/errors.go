package dunning

import (
	"errors"
)

var (
	ErrNotFound             = errors.New("dunning: not found")
	ErrInvalidState         = errors.New("dunning: invalid state transition")
	ErrNoticeFeeAmountAbove = errors.New("dunning: notice rejected, fee amount above 100000")
	ErrNoticeMissingSentAt  = errors.New("dunning: notice rejected, missing sent at")
)
