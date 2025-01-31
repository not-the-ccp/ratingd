package contract

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Contract is a contract record.
type Contract struct {
	ID         int64          `json:"id"`
	SignedAt   *time.Time     `json:"signed_at"`
	CustomerID int64          `json:"customer_id"`
	Status     ContractStatus `json:"status"`
	// ValidTo is stored in UTC.
	ValidTo           time.Time    `json:"valid_to"`
	SignedBy          string       `json:"signed_by"`
	TerminatedAt      *time.Time   `json:"terminated_at"`
	ValidFrom         time.Time    `json:"valid_from"`
	ExternalReference string       `json:"reference"`
	DiscountPercent   float64      `json:"discount_percent"`
	MinimumCharge     money.Amount `json:"minimum_charge"`
	TariffID          int64        `json:"tariff_id"`
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
