package settlement

import (
	"errors"
)

var (
	ErrNotFound                    = errors.New("settlement: not found")
	ErrInvalidState                = errors.New("settlement: invalid state transition")
	ErrSettlementDeductionsAbove   = errors.New("settlement: settlement rejected, deductions above 0")
	ErrSettlementItemDisputedIsNot = errors.New("settlement: settlement item rejected, disputed is not set")
)
