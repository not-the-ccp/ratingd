package customer

import (
	"fmt"
	"strings"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type Customer struct {
	ID              int64          `json:"id"`
	ExternalRef     string         `json:"external_ref"`
	Status          CustomerStatus `json:"status"`
	VATNumber       string         `json:"vat_number"`
	CreditLimit     money.Amount   `json:"credit_limit"`
	ActivatedAt     *time.Time     `json:"activated_at"`
	ClosedAt        *time.Time     `json:"closed_at"`
	PaymentTermDays int            `json:"payment_term_days"`
	Name            string         `json:"name"`
	BlockedAt       *time.Time     `json:"blocked_at"`
	Number          string         `json:"number"`
	Country         string         `json:"country"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type CustomerStatus string

// Statuses a customer moves through.
const (
	CustomerStatusProspect CustomerStatus = "prospect"
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusBlocked  CustomerStatus = "blocked"
	CustomerStatusClosed   CustomerStatus = "closed"
)

func (c *Customer) CanActivate() bool {
	return c.Status == CustomerStatusProspect
}

func (c *Customer) CanBlock() bool {
	return c.Status == CustomerStatusActive
}

func (s CustomerStatus) String() string {
	return string(s)
}

func (s CustomerStatus) IsFinal() bool {
	switch s {
	case CustomerStatusBlocked, CustomerStatusClosed:
		return true
	default:
		return false
	}
}

// ParseCustomerStatus converts a raw status string, ignoring case and surrounding space.
func ParseCustomerStatus(s string) (CustomerStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "prospect":
		return CustomerStatusProspect, nil
	case "active":
		return CustomerStatusActive, nil
	case "blocked":
		return CustomerStatusBlocked, nil
	case "closed":
		return CustomerStatusClosed, nil
	default:
		return "", fmt.Errorf("unknown customer status %q", s)
	}
}

func (c *Customer) HasCreditLimit() bool {
	return c.CreditLimit > 0
}
