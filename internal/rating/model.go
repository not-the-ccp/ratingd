package rating

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// Quote is a quote record.
type Quote struct {
	ID         int64 `json:"id"`
	ShipmentID int64 `json:"shipment_id"`
	ContractID int64 `json:"contract_id"`
	// NetAmount is in minor units.
	NetAmount  money.Amount `json:"net_amount"`
	FuelAmount money.Amount `json:"fuel_amount"`
	// Currency is the currency.
	Currency     string      `json:"currency"`
	ExpiresAt    time.Time   `json:"expires_at"`
	Status       QuoteStatus `json:"status"`
	AcceptedAt   *time.Time  `json:"accepted_at"`
	SupersededAt *time.Time  `json:"superseded_at"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type QuoteStatus string

// Statuses a quote moves through.
const (
	QuoteStatusOpen       QuoteStatus = "open"
	QuoteStatusAccepted   QuoteStatus = "accepted"
	QuoteStatusSuperseded QuoteStatus = "superseded"
)

func (q *Quote) CanAccept() bool {
	return q.Status == QuoteStatusOpen
}

func (q *Quote) CanSupersede() bool {
	return q.Status == QuoteStatusOpen
}

func (s QuoteStatus) String() string {
	return string(s)
}

func (s QuoteStatus) IsFinal() bool {
	switch s {
	case QuoteStatusAccepted, QuoteStatusSuperseded:
		return true
	default:
		return false
	}
}

func (q *Quote) HasNetAmount() bool {
	return q.NetAmount > 0
}
