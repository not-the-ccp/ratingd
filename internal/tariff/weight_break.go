package tariff

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// WeightBreak represents a single weight break.
type WeightBreak struct {
	ID        int64        `json:"id"`
	TariffID  int64        `json:"tariff_id"`
	FromKg    float64      `json:"from_kg"`
	ToKg      float64      `json:"to_kg"`
	Zone      string       `json:"zone"`
	RatePerKg money.Amount `json:"rate_per_kg"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func groupWeightBreaksByTariff(items []*WeightBreak) map[int64][]*WeightBreak {
	out := make(map[int64][]*WeightBreak, len(items))
	for _, item := range items {
		out[item.TariffID] = append(out[item.TariffID], item)
	}
	return out
}
