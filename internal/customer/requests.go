package customer

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateCustomerRequest struct {
	Name            string       `json:"name"`
	Number          string       `json:"number"`
	VATNumber       string       `json:"vat_number"`
	Country         string       `json:"country"`
	PaymentTermDays int          `json:"payment_term_days"`
	CreditLimit     money.Amount `json:"credit_limit"`
}

type CustomerResponse struct {
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
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// Validate checks the request before it reaches the service.
func (r *CreateCustomerRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.VATNumber == "" {
		return errors.New("vat number is required")
	}
	if r.PaymentTermDays < 0 {
		return errors.New("payment term days must not be negative")
	}
	if r.CreditLimit < 0 {
		return errors.New("credit limit must not be negative")
	}
	return nil
}

func toCustomerResponse(customer *Customer) CustomerResponse {
	return CustomerResponse{
		ID:              customer.ID,
		Name:            customer.Name,
		Number:          customer.Number,
		VATNumber:       customer.VATNumber,
		Country:         customer.Country,
		PaymentTermDays: customer.PaymentTermDays,
		CreditLimit:     customer.CreditLimit,
		Status:          customer.Status,
		ActivatedAt:     customer.ActivatedAt,
		BlockedAt:       customer.BlockedAt,
		ClosedAt:        customer.ClosedAt,
		CreatedAt:       customer.CreatedAt,
		UpdatedAt:       customer.UpdatedAt,
	}
}
