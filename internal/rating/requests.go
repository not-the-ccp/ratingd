package rating

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateQuoteRequest struct {
	ShipmentID     int64        `json:"shipment_id"`
	ContractID     int64        `json:"contract_id"`
	NetAmount      money.Amount `json:"net_amount"`
	FuelAmount     money.Amount `json:"fuel_amount"`
	Currency       string       `json:"currency"`
	ExpiresAt      time.Time    `json:"expires_at"`
	TollAmount     money.Amount `json:"toll_amount"`
	DiscountAmount money.Amount `json:"discount_amount"`
	ChargeableKg   float64      `json:"chargeable_kg"`
}

type QuoteResponse struct {
	ID             int64        `json:"id"`
	ShipmentID     int64        `json:"shipment_id"`
	ContractID     int64        `json:"contract_id"`
	NetAmount      money.Amount `json:"net_amount"`
	FuelAmount     money.Amount `json:"fuel_amount"`
	Currency       string       `json:"currency"`
	ExpiresAt      time.Time    `json:"expires_at"`
	Status         QuoteStatus  `json:"status"`
	AcceptedAt     *time.Time   `json:"accepted_at"`
	SupersededAt   *time.Time   `json:"superseded_at"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	TollAmount     money.Amount `json:"toll_amount"`
	DiscountAmount money.Amount `json:"discount_amount"`
	ChargeableKg   float64      `json:"chargeable_kg"`
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
		ID:             quote.ID,
		ShipmentID:     quote.ShipmentID,
		ContractID:     quote.ContractID,
		NetAmount:      quote.NetAmount,
		FuelAmount:     quote.FuelAmount,
		Currency:       quote.Currency,
		ExpiresAt:      quote.ExpiresAt,
		Status:         quote.Status,
		AcceptedAt:     quote.AcceptedAt,
		SupersededAt:   quote.SupersededAt,
		CreatedAt:      quote.CreatedAt,
		UpdatedAt:      quote.UpdatedAt,
		TollAmount:     quote.TollAmount,
		DiscountAmount: quote.DiscountAmount,
		ChargeableKg:   quote.ChargeableKg,
	}
}

type CreateSurchargeRequest struct {
	QuoteID     int64        `json:"quote_id"`
	Code        string       `json:"code"`
	Description string       `json:"description"`
	Amount      money.Amount `json:"amount"`
	Reason      string       `json:"reason"`
	Taxable     bool         `json:"taxable"`
}

type SurchargeResponse struct {
	ID          int64        `json:"id"`
	QuoteID     int64        `json:"quote_id"`
	Code        string       `json:"code"`
	Description string       `json:"description"`
	Amount      money.Amount `json:"amount"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Reason      string       `json:"reason"`
	Taxable     bool         `json:"taxable"`
}

func (r *CreateSurchargeRequest) Validate() error {
	if r.QuoteID == 0 {
		return errors.New("quote id is required")
	}
	if r.Code == "" {
		return errors.New("code is required")
	}
	if r.Description == "" {
		return errors.New("description is required")
	}
	if r.Amount < 0 {
		return errors.New("amount must not be negative")
	}
	return nil
}

func toSurchargeResponse(surcharge *Surcharge) SurchargeResponse {
	return SurchargeResponse{
		ID:          surcharge.ID,
		QuoteID:     surcharge.QuoteID,
		Code:        surcharge.ShortCode,
		Description: surcharge.Description,
		Amount:      surcharge.Amount,
		CreatedAt:   surcharge.CreatedAt,
		UpdatedAt:   surcharge.UpdatedAt,
		Reason:      surcharge.Reason,
		Taxable:     surcharge.Taxable,
	}
}
