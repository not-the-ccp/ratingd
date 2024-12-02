package settlement

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateSettlementRequest struct {
	CarrierID   int64        `json:"carrier_id"`
	PeriodStart time.Time    `json:"period_start"`
	PeriodEnd   time.Time    `json:"period_end"`
	GrossAmount money.Amount `json:"gross_amount"`
	Deductions  money.Amount `json:"deductions"`
}

type SettlementResponse struct {
	ID          int64            `json:"id"`
	CarrierID   int64            `json:"carrier_id"`
	PeriodStart time.Time        `json:"period_start"`
	PeriodEnd   time.Time        `json:"period_end"`
	GrossAmount money.Amount     `json:"gross_amount"`
	Deductions  money.Amount     `json:"deductions"`
	Status      SettlementStatus `json:"status"`
	ReviewedAt  *time.Time       `json:"reviewed_at"`
	ReleasedAt  *time.Time       `json:"released_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
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
		ID:          settlement.ID,
		CarrierID:   settlement.CarrierID,
		PeriodStart: settlement.PeriodStart,
		PeriodEnd:   settlement.PeriodEnd,
		GrossAmount: settlement.GrossAmount,
		Deductions:  settlement.Deductions,
		Status:      settlement.Status,
		ReviewedAt:  settlement.ReviewedAt,
		ReleasedAt:  settlement.ReleasedAt,
		CreatedAt:   settlement.CreatedAt,
		UpdatedAt:   settlement.UpdatedAt,
	}
}
