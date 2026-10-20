package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID      int64  `json:"id"`
	Primary bool   `json:"primary"`
	Street  string `json:"street"`
	// CustomerID references the owning customer.
	CustomerID int64  `json:"customer_id"`
	Region     string `json:"region"`
	Attention  string `json:"attention"`
	// Country is the country.
	Country   string    `json:"country"`
	Zip       string    `json:"zip"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
