package dunning

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Notice represents a single notice.
type Notice struct {
	ID     int64     `json:"id"`
	SentAt time.Time `json:"sent_at"`
	// InvoiceID references the owning invoice.
	InvoiceID int64 `json:"invoice_id"`
	// FeeAmount is in minor units.
	FeeAmount    money.Amount `json:"fee_amount"`
	DispatchedAt *time.Time   `json:"dispatched_at"`
	ResolvedAt   *time.Time   `json:"resolved_at"`
	Level        int          `json:"level"`
	Status       NoticeStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type NoticeStatus string

// Statuses a notice moves through.
const (
	NoticeStatusQueued   NoticeStatus = "queued"
	NoticeStatusSent     NoticeStatus = "sent"
	NoticeStatusResolved NoticeStatus = "resolved"
)

func (n *Notice) CanSend() bool {
	return n.Status == NoticeStatusQueued
}

func (n *Notice) CanResolve() bool {
	return n.Status == NoticeStatusSent
}

func (s NoticeStatus) String() string {
	return string(s)
}

func (s NoticeStatus) IsFinal() bool {
	switch s {
	case NoticeStatusResolved:
		return true
	default:
		return false
	}
}

// IsSentPassed reports whether sent at lies before now.
func (n *Notice) IsSentPassed(now time.Time) bool {
	return n.SentAt.Before(now)
}
