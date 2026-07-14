package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID                  int64      `json:"id"`
	InsuranceValidUntil time.Time  `json:"insurance_valid_until"`
	ApprovedAt          *time.Time `json:"approved_at"`
	FleetSize           int        `json:"fleet_size"`
	PaymentTermDays     int        `json:"payment_term_days"`
	ShortCode           string     `json:"code"`
	// Country is the country.
	Country     string        `json:"country"`
	VATNumber   string        `json:"vat_number"`
	Status      CarrierStatus `json:"status"`
	IBAN        string        `json:"iban"`
	SuspendedAt *time.Time    `json:"suspended_at"`
	DisplayName string        `json:"name"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type CarrierStatus string

// Statuses a carrier moves through.
const (
	CarrierStatusOnboarding CarrierStatus = "onboarding"
	CarrierStatusApproved   CarrierStatus = "approved"
	CarrierStatusSuspended  CarrierStatus = "suspended"
)

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
