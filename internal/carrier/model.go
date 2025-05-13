package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID         int64         `json:"id"`
	ApprovedAt *time.Time    `json:"approved_at"`
	Status     CarrierStatus `json:"status"`
	// Country is the country.
	Country             string     `json:"country"`
	SuspendedAt         *time.Time `json:"suspended_at"`
	PaymentTermDays     int        `json:"payment_term_days"`
	Name                string     `json:"name"`
	ShortCode           string     `json:"code"`
	VATNumber           string     `json:"vat_number"`
	InsuranceValidUntil time.Time  `json:"insurance_valid_until"`
	FleetSize           int        `json:"fleet_size"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
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
