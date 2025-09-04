package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// FuelSurcharge is a fuel surcharge record.
type FuelSurcharge struct {
	ID       int64   `json:"id"`
	Source   string  `json:"source"`
	Percent  float64 `json:"percent"`
	TariffID int64   `json:"tariff_id"`
	// DieselPrice is in minor units.
	DieselPrice money.Amount `json:"diesel_price"`
	Floor       float64      `json:"floor"`
	ValidFrom   time.Time    `json:"valid_from"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (f *FuelSurcharge) HasDieselPrice() bool {
	return f.DieselPrice > 0
}
