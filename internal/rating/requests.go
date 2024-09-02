package rating

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateQuoteRequest struct {
	ShipmentID int64        `json:"shipment_id"`
	ContractID int64        `json:"contract_id"`
	NetAmount  money.Amount `json:"net_amount"`
	FuelAmount money.Amount `json:"fuel_amount"`
	Currency   string       `json:"currency"`
	ExpiresAt  time.Time    `json:"expires_at"`
}

type QuoteResponse struct {
	ID           int64        `json:"id"`
	ShipmentID   int64        `json:"shipment_id"`
	ContractID   int64        `json:"contract_id"`
	NetAmount    money.Amount `json:"net_amount"`
	FuelAmount   money.Amount `json:"fuel_amount"`
	Currency     string       `json:"currency"`
	ExpiresAt    time.Time    `json:"expires_at"`
	Status       QuoteStatus  `json:"status"`
	AcceptedAt   *time.Time   `json:"accepted_at"`
	SupersededAt *time.Time   `json:"superseded_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// Validate checks the request before it reaches the service.
func (r *CreateQuoteRequest) Validate() error {
	if r.ShipmentID == 0 {
		return errors.New("shipment id is required")
	}
	if r.ContractID == 0 {
		return errors.New("contract id is required")
	}
	if r.NetAmount < 0 {
		return errors.New("net amount must not be negative")
	}
	if r.FuelAmount < 0 {
		return errors.New("fuel amount must not be negative")
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}

func toQuoteResponse(quote *Quote) QuoteResponse {
	return QuoteResponse{
		ID:           quote.ID,
		ShipmentID:   quote.ShipmentID,
		ContractID:   quote.ContractID,
		NetAmount:    quote.NetAmount,
		FuelAmount:   quote.FuelAmount,
		Currency:     quote.Currency,
		ExpiresAt:    quote.ExpiresAt,
		Status:       quote.Status,
		AcceptedAt:   quote.AcceptedAt,
		SupersededAt: quote.SupersededAt,
		CreatedAt:    quote.CreatedAt,
		UpdatedAt:    quote.UpdatedAt,
	}
}
