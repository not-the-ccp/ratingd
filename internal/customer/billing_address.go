package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID      int64 `json:"id"`
	Primary bool  `json:"primary"`
	// Country is the country.
	Country string `json:"country"`
	// CustomerID references the owning customer.
	CustomerID int64     `json:"customer_id"`
	Street     string    `json:"street"`
	Zip        string    `json:"zip"`
	Attention  string    `json:"attention"`
	City       string    `json:"city"`
	Region     string    `json:"region"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
