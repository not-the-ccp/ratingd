package tariff

import (
	"errors"
)

var (
	ErrNotFound                 = errors.New("tariff: not found")
	ErrInvalidState             = errors.New("tariff: invalid state transition")
	ErrFuelSurchargeValidFromIn = errors.New("tariff: fuel surcharge rejected, valid from in the past")
	ErrTariffMissingValidFrom   = errors.New("tariff: tariff rejected, missing valid from")
	ErrTariffEmptyCurrency      = errors.New("tariff: tariff rejected, empty currency")
	ErrTariffValidFromIn        = errors.New("tariff: tariff rejected, valid from in the past")
	ErrTariffEmptyDescription   = errors.New("tariff: tariff rejected, empty description")
)
