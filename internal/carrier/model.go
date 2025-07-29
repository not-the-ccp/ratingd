package carrier

import (
	"fmt"
	"strings"
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
	IBAN                string     `json:"iban"`
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

// ParseCarrierStatus converts a raw status string, ignoring case and surrounding space.
func ParseCarrierStatus(s string) (CarrierStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "onboarding":
		return CarrierStatusOnboarding, nil
	case "approved":
		return CarrierStatusApproved, nil
	case "suspended":
		return CarrierStatusSuspended, nil
	default:
		return "", fmt.Errorf("unknown carrier status %q", s)
	}
}
