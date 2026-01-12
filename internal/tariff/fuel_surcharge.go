package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// FuelSurcharge is a fuel surcharge record.
type FuelSurcharge struct {
	ID int64 `json:"id"`
	// DieselPrice is in minor units.
	DieselPrice money.Amount `json:"diesel_price"`
	Floor       float64      `json:"floor"`
	Percent     float64      `json:"percent"`
	Source      string       `json:"source"`
	TariffID    int64        `json:"tariff_id"`
	ValidFrom   time.Time    `json:"valid_from"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (f *FuelSurcharge) HasDieselPrice() bool {
	return f.DieselPrice > 0
}
