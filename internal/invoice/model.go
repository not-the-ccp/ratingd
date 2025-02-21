package invoice

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Invoice represents a single invoice.
type Invoice struct {
	ID       int64      `json:"id"`
	IssuedOn *time.Time `json:"issued_on"`
	// DueAt is stored in UTC.
	DueAt      time.Time    `json:"due_at"`
	Number     string       `json:"number"`
	CustomerID int64        `json:"customer_id"`
	VoidedAt   *time.Time   `json:"voided_at"`
	TaxAmount  money.Amount `json:"tax_amount"`
	PaidAt     *time.Time   `json:"paid_at"`
	// Currency is the currency.
	Currency string        `json:"currency"`
	Status   InvoiceStatus `json:"status"`
	// IssuedAt is stored in UTC.
	IssuedAt  time.Time    `json:"issued_at"`
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
