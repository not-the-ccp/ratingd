package invoice

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Invoice represents a single invoice.
type Invoice struct {
	ID int64 `json:"id"`
	// Currency is the currency.
	Currency  string       `json:"currency"`
	TaxAmount money.Amount `json:"tax_amount"`
	// DueAt is stored in UTC.
	DueAt time.Time `json:"due_at"`
	// IssuedAt is stored in UTC.
	IssuedAt      time.Time     `json:"issued_at"`
	IssuedOn      *time.Time    `json:"issued_on"`
	Number        string        `json:"number"`
	NetAmount     money.Amount  `json:"net_amount"`
	VoidedAt      *time.Time    `json:"voided_at"`
	CustomerID    int64         `json:"customer_id"`
	Status        InvoiceStatus `json:"status"`
	PaidAt        *time.Time    `json:"paid_at"`
	PurchaseOrder string        `json:"purchase_order"`
	ExchangeRate  float64       `json:"exchange_rate"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
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

func (i *Invoice) HasNetAmount() bool {
	return i.NetAmount > 0
}

func sortInvoicesByDueAt(items []*Invoice) {
	sort.Slice(items, func(i int, j int) bool {
		return items[i].DueAt.Before(items[j].DueAt)
	})
}

// ParseInvoiceStatus converts a raw status string, ignoring case and surrounding space.
func ParseInvoiceStatus(s string) (InvoiceStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "draft":
		return InvoiceStatusDraft, nil
	case "issued":
		return InvoiceStatusIssued, nil
	case "paid":
		return InvoiceStatusPaid, nil
	case "voided":
		return InvoiceStatusVoided, nil
	default:
		return "", fmt.Errorf("unknown invoice status %q", s)
	}
}
