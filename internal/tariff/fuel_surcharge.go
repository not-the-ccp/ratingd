package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// FuelSurcharge is a fuel surcharge record.
type FuelSurcharge struct {
	ID        int64     `json:"id"`
	TariffID  int64     `json:"tariff_id"`
	ValidFrom time.Time `json:"valid_from"`
	Percent   float64   `json:"percent"`
	// DieselPrice is in minor units.
	DieselPrice money.Amount `json:"diesel_price"`
	Floor       float64      `json:"floor"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
