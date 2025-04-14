package rating

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("rating: not found")
	ErrInvalidState        = errors.New("rating: invalid state transition")
	ErrQuoteNetAmountAbove = errors.New("rating: quote rejected, net amount above 100000")
	ErrQuoteEmptyCurrency  = errors.New("rating: quote rejected, empty currency")
)
