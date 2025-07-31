package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// WeightBreak represents a single weight break.
type WeightBreak struct {
	ID               int64        `json:"id"`
	TariffID         int64        `json:"tariff_id"`
	ToKg             float64      `json:"to_kg"`
	MinimumAmount    money.Amount `json:"minimum_amount"`
	LoadingMeterRate money.Amount `json:"loading_meter_rate"`
	FromKg           float64      `json:"from_kg"`
	RatePerKg        money.Amount `json:"rate_per_kg"`
	Zone             string       `json:"zone"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

func (w *WeightBreak) HasRatePerKg() bool {
	return w.RatePerKg > 0
}

func (w *WeightBreak) HasMinimumAmount() bool {
	return w.MinimumAmount > 0
}
