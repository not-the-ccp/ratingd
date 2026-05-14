package carrier

import (
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID                  int64         `json:"id"`
	ApprovedAt          *time.Time    `json:"approved_at"`
	VATNumber           string        `json:"vat_number"`
	Status              CarrierStatus `json:"status"`
	PaymentTermDays     int           `json:"payment_term_days"`
	SuspendedAt         *time.Time    `json:"suspended_at"`
	FleetSize           int           `json:"fleet_size"`
	DisplayName         string        `json:"name"`
	InsuranceValidUntil time.Time     `json:"insurance_valid_until"`
	ShortCode           string        `json:"code"`
	// Country is the country.
	Country   string    `json:"country"`
	IBAN      string    `json:"iban"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
