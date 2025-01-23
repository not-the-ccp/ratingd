package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID         int64      `json:"id"`
	ApprovedAt *time.Time `json:"approved_at"`
	// Country is the country.
	Country         string        `json:"country"`
	ShortCode       string        `json:"code"`
	Name            string        `json:"name"`
	VATNumber       string        `json:"vat_number"`
	SuspendedAt     *time.Time    `json:"suspended_at"`
	Status          CarrierStatus `json:"status"`
	PaymentTermDays int           `json:"payment_term_days"`
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
