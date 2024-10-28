package shipment

import (
	"errors"
)

var (
	ErrNotFound                = errors.New("shipment: not found")
	ErrInvalidState            = errors.New("shipment: invalid state transition")
	ErrShipmentOriginCountryCh = errors.New("shipment: shipment rejected, origin country CH")
)
