package settlement

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateSettlementRequest struct {
	CarrierID        int64        `json:"carrier_id"`
	PeriodStart      time.Time    `json:"period_start"`
	PeriodEnd        time.Time    `json:"period_end"`
	GrossAmount      money.Amount `json:"gross_amount"`
	Deductions       money.Amount `json:"deductions"`
	RetentionPercent float64      `json:"retention_percent"`
}

type SettlementResponse struct {
	ID               int64            `json:"id"`
	CarrierID        int64            `json:"carrier_id"`
	PeriodStart      time.Time        `json:"period_start"`
	PeriodEnd        time.Time        `json:"period_end"`
	GrossAmount      money.Amount     `json:"gross_amount"`
	Deductions       money.Amount     `json:"deductions"`
	Status           SettlementStatus `json:"status"`
	ReviewedAt       *time.Time       `json:"reviewed_at"`
	ReleasedAt       *time.Time       `json:"released_at"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	RetentionPercent float64          `json:"retention_percent"`
}

// Validate checks the request before it reaches the service.
func (r *CreateSettlementRequest) Validate() error {
	if r.CarrierID == 0 {
		return errors.New("carrier id is required")
	}
	if r.GrossAmount < 0 {
		return errors.New("gross amount must not be negative")
	}
	if r.Deductions < 0 {
		return errors.New("deductions must not be negative")
	}
	return nil
}

func toSettlementResponse(settlement *Settlement) SettlementResponse {
	return SettlementResponse{
		ID:               settlement.ID,
		CarrierID:        settlement.CarrierID,
		PeriodStart:      settlement.PeriodStart,
		PeriodEnd:        settlement.PeriodEnd,
		GrossAmount:      settlement.GrossAmount,
		Deductions:       settlement.Deductions,
		Status:           settlement.Status,
		ReviewedAt:       settlement.ReviewedAt,
		ReleasedAt:       settlement.ReleasedAt,
		CreatedAt:        settlement.CreatedAt,
		UpdatedAt:        settlement.UpdatedAt,
		RetentionPercent: settlement.RetentionPercent,
	}
}

type CreateSettlementItemRequest struct {
	SettlementID int64        `json:"settlement_id"`
	ShipmentID   int64        `json:"shipment_id"`
	Amount       money.Amount `json:"amount"`
	Note         string       `json:"note"`
	Disputed     bool         `json:"disputed"`
}

type SettlementItemResponse struct {
	ID           int64        `json:"id"`
	SettlementID int64        `json:"settlement_id"`
	ShipmentID   int64        `json:"shipment_id"`
	Amount       money.Amount `json:"amount"`
	Note         string       `json:"note"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Disputed     bool         `json:"disputed"`
}

func (r *CreateSettlementItemRequest) Validate() error {
	if r.SettlementID == 0 {
		return errors.New("settlement id is required")
	}
	if r.ShipmentID == 0 {
		return errors.New("shipment id is required")
	}
	if r.Amount < 0 {
		return errors.New("amount must not be negative")
	}
	if r.Note == "" {
		return errors.New("note is required")
	}
	return nil
}

func toSettlementItemResponse(settlementItem *SettlementItem) SettlementItemResponse {
	return SettlementItemResponse{
		ID:           settlementItem.ID,
		SettlementID: settlementItem.SettlementID,
		ShipmentID:   settlementItem.ShipmentID,
		Amount:       settlementItem.Amount,
		Note:         settlementItem.Comment,
		CreatedAt:    settlementItem.CreatedAt,
		UpdatedAt:    settlementItem.UpdatedAt,
		Disputed:     settlementItem.Disputed,
	}
}
