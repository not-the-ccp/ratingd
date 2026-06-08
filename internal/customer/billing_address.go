package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID        int64  `json:"id"`
	Attention string `json:"attention"`
	// CustomerID references the owning customer.
	CustomerID int64  `json:"customer_id"`
	City       string `json:"city"`
	Zip        string `json:"zip"`
	// Country is the country.
	Country   string    `json:"country"`
	Region    string    `json:"region"`
	Primary   bool      `json:"primary"`
	Street    string    `json:"street"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
