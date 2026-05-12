package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID      int64 `json:"id"`
	Primary bool  `json:"primary"`
	// CustomerID references the owning customer.
	CustomerID int64  `json:"customer_id"`
	Zip        string `json:"zip"`
	Region     string `json:"region"`
	Street     string `json:"street"`
	// Country is the country.
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Attention string    `json:"attention"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
