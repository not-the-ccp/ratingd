package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID int64 `json:"id"`
	// Country is the country.
	Country   string `json:"country"`
	Region    string `json:"region"`
	Street    string `json:"street"`
	Zip       string `json:"zip"`
	City      string `json:"city"`
	Primary   bool   `json:"primary"`
	Attention string `json:"attention"`
	// CustomerID references the owning customer.
	CustomerID int64     `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
