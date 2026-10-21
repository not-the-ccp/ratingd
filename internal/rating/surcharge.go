package rating

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Surcharge represents a single surcharge.
type Surcharge struct {
	ID          int64        `json:"id"`
	Description string       `json:"description"`
	Reason      string       `json:"reason"`
	Taxable     bool         `json:"taxable"`
	QuoteID     int64        `json:"quote_id"`
	ShortCode   string       `json:"code"`
	Amount      money.Amount `json:"amount"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (s *Surcharge) HasAmount() bool {
	return s.Amount > 0
}
