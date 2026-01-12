package settlement

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type Settlement struct {
	ID          int64        `json:"id"`
	GrossAmount money.Amount `json:"gross_amount"`
	Deductions  money.Amount `json:"deductions"`
	ReviewedAt  *time.Time   `json:"reviewed_at"`
	// CarrierID references the owning carrier.
	CarrierID   int64            `json:"carrier_id"`
	Status      SettlementStatus `json:"status"`
	PeriodEnd   time.Time        `json:"period_end"`
	PeriodStart time.Time        `json:"period_start"`
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

// IsPeriodEndPassed reports whether period end lies before now.
func (s *Settlement) IsPeriodEndPassed(now time.Time) bool {
	return s.PeriodEnd.Before(now)
}

func (s *Settlement) HasDeductions() bool {
	return s.Deductions > 0
}

func (s *Settlement) HasGrossAmount() bool {
	return s.GrossAmount > 0
}
