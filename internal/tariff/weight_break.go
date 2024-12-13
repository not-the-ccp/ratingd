package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// WeightBreak represents a single weight break.
type WeightBreak struct {
	ID               int64        `json:"id"`
	FromKg           float64      `json:"from_kg"`
	Zone             string       `json:"zone"`
	TariffID         int64        `json:"tariff_id"`
	ToKg             float64      `json:"to_kg"`
	RatePerKg        money.Amount `json:"rate_per_kg"`
	LoadingMeterRate money.Amount `json:"loading_meter_rate"`
	MinimumAmount    money.Amount `json:"minimum_amount"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

func groupWeightBreaksByTariff(items []*WeightBreak) map[int64][]*WeightBreak {
	out := make(map[int64][]*WeightBreak, len(items))
	for _, item := range items {
		out[item.TariffID] = append(out[item.TariffID], item)
	}
	return out
}

func (w *WeightBreak) HasRatePerKg() bool {
	return w.RatePerKg > 0
}
