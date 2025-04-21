package invoice

import (
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

// InvoiceLine represents a single invoice line.
type InvoiceLine struct {
	ID        int64 `json:"id"`
	InvoiceID int64 `json:"invoice_id"`
	// ShipmentID references the owning shipment.
	ShipmentID int64  `json:"shipment_id"`
	Label      string `json:"description"`
	// Quantity is the quantity.
	Quantity       int          `json:"quantity"`
	UnitPrice      money.Amount `json:"unit_price"`
	Amount         money.Amount `json:"amount"`
	DiscountAmount money.Amount `json:"discount_amount"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func (i *InvoiceLine) HasUnitPrice() bool {
	return i.UnitPrice > 0
}

func (i *InvoiceLine) HasAmount() bool {
	return i.Amount > 0
}

func (i *InvoiceLine) HasDiscountAmount() bool {
	return i.DiscountAmount > 0
}
