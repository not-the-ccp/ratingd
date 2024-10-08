package customer

import (
	"time"
)

// BillingAddress holds the persisted state of a billing address.
type BillingAddress struct {
	ID int64 `json:"id"`
	// CustomerID references the owning customer.
	CustomerID int64  `json:"customer_id"`
	City       string `json:"city"`
	Attention  string `json:"attention"`
	Zip        string `json:"zip"`
	Street     string `json:"street"`
	Primary    bool   `json:"primary"`
	// Country is the country.
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
