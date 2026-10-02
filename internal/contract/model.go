package contract

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Contract is a contract record.
type Contract struct {
	ID         int64          `json:"id"`
	ValidFrom  time.Time      `json:"valid_from"`
	CustomerID int64          `json:"customer_id"`
	Status     ContractStatus `json:"status"`
	// ValidTo is stored in UTC.
	ValidTo           time.Time    `json:"valid_to"`
	DiscountPercent   float64      `json:"discount_percent"`
	ExternalReference string       `json:"reference"`
	TerminatedAt      *time.Time   `json:"terminated_at"`
	SignedBy          string       `json:"signed_by"`
	SignedAt          *time.Time   `json:"signed_at"`
	TariffID          int64        `json:"tariff_id"`
	MinimumCharge     money.Amount `json:"minimum_charge"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type ContractStatus string

// Statuses a contract moves through.
const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusSigned     ContractStatus = "signed"
	ContractStatusTerminated ContractStatus = "terminated"
	ContractStatusExpired    ContractStatus = "expired"
)

func (c *Contract) CanSign() bool {
	return c.Status == ContractStatusDraft
}

func (c *Contract) CanTerminate() bool {
	return c.Status == ContractStatusSigned
}

func (s ContractStatus) String() string {
	return string(s)
}

func (s ContractStatus) IsFinal() bool {
	switch s {
	case ContractStatusTerminated, ContractStatusExpired:
		return true
	default:
		return false
	}
}

// IsValidToPassed reports whether valid to lies before now.
func (c *Contract) IsValidToPassed(now time.Time) bool {
	return c.ValidTo.Before(now)
}

func (c *Contract) HasMinimumCharge() bool {
	return c.MinimumCharge > 0
}
