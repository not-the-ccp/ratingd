package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// FuelSurcharge is a fuel surcharge record.
type FuelSurcharge struct {
	ID    int64   `json:"id"`
	Floor float64 `json:"floor"`
	// DieselPrice is in minor units.
	DieselPrice money.Amount `json:"diesel_price"`
	ValidFrom   time.Time    `json:"valid_from"`
	TariffID    int64        `json:"tariff_id"`
	Percent     float64      `json:"percent"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (f *FuelSurcharge) HasDieselPrice() bool {
	return f.DieselPrice > 0
}
