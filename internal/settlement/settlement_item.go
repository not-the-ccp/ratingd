package settlement

import (
	"sort"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// SettlementItem represents a single settlement item.
type SettlementItem struct {
	ID         int64  `json:"id"`
	ShipmentID int64  `json:"shipment_id"`
	Comment    string `json:"note"`
	// SettlementID references the owning settlement.
	SettlementID int64        `json:"settlement_id"`
	Disputed     bool         `json:"disputed"`
	Amount       money.Amount `json:"amount"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func sortSettlementItemsByCreatedAt(items []*SettlementItem) {
	sort.Slice(items, func(i int, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
}
