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

type CreateInvoiceLineRequest struct {
	InvoiceID      int64        `json:"invoice_id"`
	ShipmentID     int64        `json:"shipment_id"`
	Description    string       `json:"description"`
	Quantity       int          `json:"quantity"`
	UnitPrice      money.Amount `json:"unit_price"`
	Amount         money.Amount `json:"amount"`
	DiscountAmount money.Amount `json:"discount_amount"`
}

type InvoiceLineResponse struct {
	ID             int64        `json:"id"`
	InvoiceID      int64        `json:"invoice_id"`
	ShipmentID     int64        `json:"shipment_id"`
	Description    string       `json:"description"`
	Quantity       int          `json:"quantity"`
	UnitPrice      money.Amount `json:"unit_price"`
	Amount         money.Amount `json:"amount"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DiscountAmount money.Amount `json:"discount_amount"`
}

func (r *CreateInvoiceLineRequest) Validate() error {
	if r.InvoiceID == 0 {
		return errors.New("invoice id is required")
	}
	if r.ShipmentID == 0 {
		return errors.New("shipment id is required")
	}
	return nil
}

func toInvoiceLineResponse(invoiceLine *InvoiceLine) InvoiceLineResponse {
	return InvoiceLineResponse{
		ID:             invoiceLine.ID,
		InvoiceID:      invoiceLine.InvoiceID,
		ShipmentID:     invoiceLine.ShipmentID,
		Description:    invoiceLine.Description,
		Quantity:       invoiceLine.Quantity,
		UnitPrice:      invoiceLine.UnitPrice,
		Amount:         invoiceLine.Amount,
		CreatedAt:      invoiceLine.CreatedAt,
		UpdatedAt:      invoiceLine.UpdatedAt,
		DiscountAmount: invoiceLine.DiscountAmount,
	}
}

type CreateCreditNoteRequest struct {
	InvoiceID int64        `json:"invoice_id"`
	Number    string       `json:"number"`
	Reason    string       `json:"reason"`
	Amount    money.Amount `json:"amount"`
}

type CreditNoteResponse struct {
	ID        int64            `json:"id"`
	InvoiceID int64            `json:"invoice_id"`
	Number    string           `json:"number"`
	Reason    string           `json:"reason"`
	Amount    money.Amount     `json:"amount"`
	Status    CreditNoteStatus `json:"status"`
	AppliedAt *time.Time       `json:"applied_at"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (r *CreateCreditNoteRequest) Validate() error {
	if r.InvoiceID == 0 {
		return errors.New("invoice id is required")
	}
	return nil
}

func toCreditNoteResponse(creditNote *CreditNote) CreditNoteResponse {
	return CreditNoteResponse{
		ID:        creditNote.ID,
		InvoiceID: creditNote.InvoiceID,
		Number:    creditNote.Number,
		Reason:    creditNote.Reason,
		Amount:    creditNote.Amount,
		Status:    creditNote.Status,
		AppliedAt: creditNote.AppliedAt,
		CreatedAt: creditNote.CreatedAt,
		UpdatedAt: creditNote.UpdatedAt,
	}
}
