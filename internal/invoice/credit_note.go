package invoice

import (
	"fmt"
	"strings"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreditNote struct {
	ID int64 `json:"id"`
	// Reason is the reason.
	Reason    string           `json:"reason"`
	Amount    money.Amount     `json:"amount"`
	Status    CreditNoteStatus `json:"status"`
	Number    string           `json:"number"`
	AppliedAt *time.Time       `json:"applied_at"`
	// InvoiceID references the owning invoice.
	InvoiceID int64     `json:"invoice_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

func (c *CreditNote) HasAmount() bool {
	return c.Amount > 0
}

// ParseCreditNoteStatus converts a raw status string, ignoring case and surrounding space.
func ParseCreditNoteStatus(s string) (CreditNoteStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "open":
		return CreditNoteStatusOpen, nil
	case "applied":
		return CreditNoteStatusApplied, nil
	default:
		return "", fmt.Errorf("unknown credit note status %q", s)
	}
}
