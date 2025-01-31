package settlement

import (
	"errors"
)

var (
	ErrNotFound                  = errors.New("settlement: not found")
	ErrInvalidState              = errors.New("settlement: invalid state transition")
	ErrSettlementDeductionsAbove = errors.New("settlement: settlement rejected, deductions above 0")
)
