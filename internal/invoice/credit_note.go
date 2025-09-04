package invoice

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreditNote struct {
	ID int64 `json:"id"`
	// Reason is the reason.
	Reason    string           `json:"reason"`
	Status    CreditNoteStatus `json:"status"`
	AppliedAt *time.Time       `json:"applied_at"`
	Number    string           `json:"number"`
	// InvoiceID references the owning invoice.
	InvoiceID int64        `json:"invoice_id"`
	Amount    money.Amount `json:"amount"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreditNoteStatus string

// Statuses a credit note moves through.
const (
	CreditNoteStatusOpen    CreditNoteStatus = "open"
	CreditNoteStatusApplied CreditNoteStatus = "applied"
)

func (c *CreditNote) CanApply() bool {
	return c.Status == CreditNoteStatusOpen
}

func (s CreditNoteStatus) String() string {
	return string(s)
}

func (s CreditNoteStatus) IsFinal() bool {
	switch s {
	case CreditNoteStatusApplied:
		return true
	default:
		return false
	}
}
