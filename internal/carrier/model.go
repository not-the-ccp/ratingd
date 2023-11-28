package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID int64 `json:"id"`
	// Country is the country.
	Country         string        `json:"country"`
	ShortCode       string        `json:"code"`
	VATNumber       string        `json:"vat_number"`
	ApprovedAt      *time.Time    `json:"approved_at"`
	Name            string        `json:"name"`
	Status          CarrierStatus `json:"status"`
	PaymentTermDays int           `json:"payment_term_days"`
	SuspendedAt     *time.Time    `json:"suspended_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type CarrierStatus string

// Statuses a carrier moves through.
const (
	CarrierStatusOnboarding CarrierStatus = "onboarding"
	CarrierStatusApproved   CarrierStatus = "approved"
	CarrierStatusSuspended  CarrierStatus = "suspended"
)

func (c *Carrier) CanApprove() bool {
	return c.Status == CarrierStatusOnboarding
}

func (c *Carrier) CanSuspend() bool {
	return c.Status == CarrierStatusApproved
}

func (s CarrierStatus) String() string {
	return string(s)
}

func (s CarrierStatus) IsFinal() bool {
	switch s {
	case CarrierStatusSuspended:
		return true
	default:
		return false
	}
}
