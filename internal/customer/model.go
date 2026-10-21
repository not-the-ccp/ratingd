package customer

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type Customer struct {
	ID              int64          `json:"id"`
	ActivatedAt     *time.Time     `json:"activated_at"`
	ExternalRef     string         `json:"external_ref"`
	VATNumber       string         `json:"vat_number"`
	DisplayName     string         `json:"name"`
	ClosedAt        *time.Time     `json:"closed_at"`
	BlockedAt       *time.Time     `json:"blocked_at"`
	DocumentNumber  string         `json:"number"`
	Status          CustomerStatus `json:"status"`
	Country         string         `json:"country"`
	IBAN            string         `json:"iban"`
	PaymentTermDays int            `json:"payment_term_days"`
	CreditLimit     money.Amount   `json:"credit_limit"`
	BillingEmail    string         `json:"billing_email"`
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

func (c *Customer) HasCreditLimit() bool {
	return c.CreditLimit > 0
}
