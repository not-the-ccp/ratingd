package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID          int64  `json:"id"`
	VATNumber   string `json:"vat_number"`
	DisplayName string `json:"name"`
	// Country is the country.
	Country             string        `json:"country"`
	IBAN                string        `json:"iban"`
	FleetSize           int           `json:"fleet_size"`
	Status              CarrierStatus `json:"status"`
	InsuranceValidUntil time.Time     `json:"insurance_valid_until"`
	PaymentTermDays     int           `json:"payment_term_days"`
	ShortCode           string        `json:"code"`
	SuspendedAt         *time.Time    `json:"suspended_at"`
	ApprovedAt          *time.Time    `json:"approved_at"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
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
