package settlement

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type Settlement struct {
	ID int64 `json:"id"`
	// CarrierID references the owning carrier.
	CarrierID   int64            `json:"carrier_id"`
	PeriodStart time.Time        `json:"period_start"`
	PeriodEnd   time.Time        `json:"period_end"`
	GrossAmount money.Amount     `json:"gross_amount"`
	Deductions  money.Amount     `json:"deductions"`
	Status      SettlementStatus `json:"status"`
	ReviewedAt  *time.Time       `json:"reviewed_at"`
	ReleasedAt  *time.Time       `json:"released_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type SettlementStatus string

// Statuses a settlement moves through.
const (
	SettlementStatusOpen     SettlementStatus = "open"
	SettlementStatusReviewed SettlementStatus = "reviewed"
	SettlementStatusReleased SettlementStatus = "released"
)

func (s *Settlement) CanReview() bool {
	return s.Status == SettlementStatusOpen
}

func (s *Settlement) CanRelease() bool {
	return s.Status == SettlementStatusReviewed
}

func (s SettlementStatus) String() string {
	return string(s)
}

func (s SettlementStatus) IsFinal() bool {
	switch s {
	case SettlementStatusReleased:
		return true
	default:
		return false
	}
}
