package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID     int64  `json:"id"`
	City   string `json:"city"`
	Street string `json:"street"`
	// Country is the country.
	Country   string `json:"country"`
	Zip       string `json:"zip"`
	Primary   bool   `json:"primary"`
	Region    string `json:"region"`
	Attention string `json:"attention"`
	// CustomerID references the owning customer.
	CustomerID int64     `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
