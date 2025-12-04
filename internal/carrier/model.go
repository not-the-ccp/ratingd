package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID                  int64      `json:"id"`
	SuspendedAt         *time.Time `json:"suspended_at"`
	VATNumber           string     `json:"vat_number"`
	IBAN                string     `json:"iban"`
	ShortCode           string     `json:"code"`
	DisplayName         string     `json:"name"`
	PaymentTermDays     int        `json:"payment_term_days"`
	InsuranceValidUntil time.Time  `json:"insurance_valid_until"`
	FleetSize           int        `json:"fleet_size"`
	ApprovedAt          *time.Time `json:"approved_at"`
	// Country is the country.
	Country   string        `json:"country"`
	Status    CarrierStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
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
