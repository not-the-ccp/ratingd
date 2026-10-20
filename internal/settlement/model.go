package settlement

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type Settlement struct {
	ID               int64            `json:"id"`
	PaymentRun       string           `json:"payment_run"`
	Deductions       money.Amount     `json:"deductions"`
	Status           SettlementStatus `json:"status"`
	ReleasedAt       *time.Time       `json:"released_at"`
	RetentionPercent float64          `json:"retention_percent"`
	Remark           string           `json:"remark"`
	GrossAmount      money.Amount     `json:"gross_amount"`
	// CarrierID references the owning carrier.
	CarrierID   int64      `json:"carrier_id"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	PeriodStart time.Time  `json:"period_start"`
	PeriodEnd   time.Time  `json:"period_end"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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

func ParseSettlementStatus(s string) (SettlementStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "open":
		return SettlementStatusOpen, nil
	case "reviewed":
		return SettlementStatusReviewed, nil
	case "released":
		return SettlementStatusReleased, nil
	default:
		return "", fmt.Errorf("unknown settlement status %q", s)
	}
}

func sortSettlementsByPeriodEnd(items []*Settlement) {
	sort.Slice(items, func(i int, j int) bool {
		return items[i].PeriodEnd.Before(items[j].PeriodEnd)
	})
}
