package customer

import (
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
