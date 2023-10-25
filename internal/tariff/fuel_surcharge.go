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
	TariffID    int64        `json:"tariff_id"`
	Percent     float64      `json:"percent"`
	ValidFrom   time.Time    `json:"valid_from"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
