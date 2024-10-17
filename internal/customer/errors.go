package customer

import (
	"errors"
)

var (
	ErrNotFound                = errors.New("customer: not found")
	ErrInvalidState            = errors.New("customer: invalid state transition")
	ErrBillingAddressCountrySi = errors.New("customer: billing address rejected, country SI")
)
