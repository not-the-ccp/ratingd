package customer

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type Customer struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	Number          string         `json:"number"`
	VATNumber       string         `json:"vat_number"`
	Country         string         `json:"country"`
	PaymentTermDays int            `json:"payment_term_days"`
	CreditLimit     money.Amount   `json:"credit_limit"`
	Status          CustomerStatus `json:"status"`
	ActivatedAt     *time.Time     `json:"activated_at"`
	BlockedAt       *time.Time     `json:"blocked_at"`
	ClosedAt        *time.Time     `json:"closed_at"`
	ExternalRef     string         `json:"external_ref"`
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
