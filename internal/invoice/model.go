package invoice

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Invoice represents a single invoice.
type Invoice struct {
	ID         int64        `json:"id"`
	CustomerID int64        `json:"customer_id"`
	TaxAmount  money.Amount `json:"tax_amount"`
	// Currency is the currency.
	Currency string        `json:"currency"`
	Status   InvoiceStatus `json:"status"`
	VoidedAt *time.Time    `json:"voided_at"`
	Number   string        `json:"number"`
	PaidAt   *time.Time    `json:"paid_at"`
	// DueAt is stored in UTC.
	DueAt time.Time `json:"due_at"`
	// IssuedAt is stored in UTC.
	IssuedAt  time.Time    `json:"issued_at"`
	IssuedOn  *time.Time   `json:"issued_on"`
	NetAmount money.Amount `json:"net_amount"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type InvoiceStatus string

// Statuses a invoice moves through.
const (
	InvoiceStatusDraft  InvoiceStatus = "draft"
	InvoiceStatusIssued InvoiceStatus = "issued"
	InvoiceStatusPaid   InvoiceStatus = "paid"
	InvoiceStatusVoided InvoiceStatus = "voided"
)

func (i *Invoice) CanIssue() bool {
	return i.Status == InvoiceStatusDraft
}

func (i *Invoice) CanMarkPaid() bool {
	return i.Status == InvoiceStatusIssued
}

func (s InvoiceStatus) String() string {
	return string(s)
}

func (s InvoiceStatus) IsFinal() bool {
	switch s {
	case InvoiceStatusPaid, InvoiceStatusVoided:
		return true
	default:
		return false
	}
}

// IsIssuedPassed reports whether issued at lies before now.
func (i *Invoice) IsIssuedPassed(now time.Time) bool {
	return i.IssuedAt.Before(now)
}

func (i *Invoice) IsDuePassed(now time.Time) bool {
	return i.DueAt.Before(now)
}

func (i *Invoice) HasTaxAmount() bool {
	return i.TaxAmount > 0
}
