package invoice

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateInvoiceRequest struct {
	CustomerID int64        `json:"customer_id"`
	Number     string       `json:"number"`
	Currency   string       `json:"currency"`
	IssuedAt   time.Time    `json:"issued_at"`
	DueAt      time.Time    `json:"due_at"`
	NetAmount  money.Amount `json:"net_amount"`
	TaxAmount  money.Amount `json:"tax_amount"`
}

type InvoiceResponse struct {
	ID         int64         `json:"id"`
	CustomerID int64         `json:"customer_id"`
	Number     string        `json:"number"`
	Currency   string        `json:"currency"`
	IssuedAt   time.Time     `json:"issued_at"`
	DueAt      time.Time     `json:"due_at"`
	NetAmount  money.Amount  `json:"net_amount"`
	TaxAmount  money.Amount  `json:"tax_amount"`
	Status     InvoiceStatus `json:"status"`
	IssuedOn   *time.Time    `json:"issued_on"`
	PaidAt     *time.Time    `json:"paid_at"`
	VoidedAt   *time.Time    `json:"voided_at"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// Validate checks the request before it reaches the service.
func (r *CreateInvoiceRequest) Validate() error {
	if r.CustomerID == 0 {
		return errors.New("customer id is required")
	}
	if r.TaxAmount < 0 {
		return errors.New("tax amount must not be negative")
	}
	return nil
}

func toInvoiceResponse(invoice *Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID:         invoice.ID,
		CustomerID: invoice.CustomerID,
		Number:     invoice.Number,
		Currency:   invoice.Currency,
		IssuedAt:   invoice.IssuedAt,
		DueAt:      invoice.DueAt,
		NetAmount:  invoice.NetAmount,
		TaxAmount:  invoice.TaxAmount,
		Status:     invoice.Status,
		IssuedOn:   invoice.IssuedOn,
		PaidAt:     invoice.PaidAt,
		VoidedAt:   invoice.VoidedAt,
		CreatedAt:  invoice.CreatedAt,
		UpdatedAt:  invoice.UpdatedAt,
	}
}
