package carrier

import (
	"fmt"
	"strings"
	"time"
)

// Carrier represents a single carrier.
type Carrier struct {
	ID                  int64         `json:"id"`
	SuspendedAt         *time.Time    `json:"suspended_at"`
	IBAN                string        `json:"iban"`
	ShortCode           string        `json:"code"`
	PaymentTermDays     int           `json:"payment_term_days"`
	Status              CarrierStatus `json:"status"`
	InsuranceValidUntil time.Time     `json:"insurance_valid_until"`
	Name                string        `json:"name"`
	ApprovedAt          *time.Time    `json:"approved_at"`
	// Country is the country.
	Country   string    `json:"country"`
	VATNumber string    `json:"vat_number"`
	FleetSize int       `json:"fleet_size"`
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
