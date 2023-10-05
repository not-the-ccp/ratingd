package tariff

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateTariffRequest struct {
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
}

type TariffResponse struct {
	ID          int64        `json:"id"`
	Code        string       `json:"code"`
	Description string       `json:"description"`
	Currency    string       `json:"currency"`
	ValidFrom   time.Time    `json:"valid_from"`
	ValidTo     time.Time    `json:"valid_to"`
	Status      TariffStatus `json:"status"`
	PublishedAt *time.Time   `json:"published_at"`
	RetiredAt   *time.Time   `json:"retired_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (r *CreateTariffRequest) Validate() error {
	if r.Description == "" {
		return errors.New("description is required")
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}

func toTariffResponse(tariff *Tariff) TariffResponse {
	return TariffResponse{
		ID:          tariff.ID,
		Code:        tariff.Code,
		Description: tariff.Description,
		Currency:    tariff.Currency,
		ValidFrom:   tariff.ValidFrom,
		ValidTo:     tariff.ValidTo,
		Status:      tariff.Status,
		PublishedAt: tariff.PublishedAt,
		RetiredAt:   tariff.RetiredAt,
		CreatedAt:   tariff.CreatedAt,
		UpdatedAt:   tariff.UpdatedAt,
	}
}

type CreateWeightBreakRequest struct {
	TariffID  int64        `json:"tariff_id"`
	FromKg    float64      `json:"from_kg"`
	ToKg      float64      `json:"to_kg"`
	Zone      string       `json:"zone"`
	RatePerKg money.Amount `json:"rate_per_kg"`
}

type WeightBreakResponse struct {
	ID        int64        `json:"id"`
	TariffID  int64        `json:"tariff_id"`
	FromKg    float64      `json:"from_kg"`
	ToKg      float64      `json:"to_kg"`
	Zone      string       `json:"zone"`
	RatePerKg money.Amount `json:"rate_per_kg"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (r *CreateWeightBreakRequest) Validate() error {
	if r.TariffID == 0 {
		return errors.New("tariff id is required")
	}
	if r.FromKg < 0 {
		return errors.New("from kg must not be negative")
	}
	if r.ToKg < 0 {
		return errors.New("to kg must not be negative")
	}
	if r.Zone == "" {
		return errors.New("zone is required")
	}
	if r.RatePerKg < 0 {
		return errors.New("rate per kg must not be negative")
	}
	return nil
}

func toWeightBreakResponse(weightBreak *WeightBreak) WeightBreakResponse {
	return WeightBreakResponse{
		ID:        weightBreak.ID,
		TariffID:  weightBreak.TariffID,
		FromKg:    weightBreak.FromKg,
		ToKg:      weightBreak.ToKg,
		Zone:      weightBreak.Zone,
		RatePerKg: weightBreak.RatePerKg,
		CreatedAt: weightBreak.CreatedAt,
		UpdatedAt: weightBreak.UpdatedAt,
	}
}
