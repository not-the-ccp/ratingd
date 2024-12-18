package settlement

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// SettlementItem represents a single settlement item.
type SettlementItem struct {
	ID int64 `json:"id"`
	// SettlementID references the owning settlement.
	SettlementID int64        `json:"settlement_id"`
	ShipmentID   int64        `json:"shipment_id"`
	Amount       money.Amount `json:"amount"`
	Note         string       `json:"note"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
