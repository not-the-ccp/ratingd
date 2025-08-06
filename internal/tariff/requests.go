package tariff

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateTariffRequest struct {
	Code         string    `json:"code,omitempty"`
	Description  string    `json:"description,omitempty"`
	Currency     string    `json:"currency,omitempty"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
	RoundingMode string    `json:"rounding_mode,omitempty"`
}

type TariffResponse struct {
	ID           int64        `json:"id"`
	Code         string       `json:"code,omitempty"`
	Description  string       `json:"description,omitempty"`
	Currency     string       `json:"currency,omitempty"`
	ValidFrom    time.Time    `json:"valid_from"`
	ValidTo      time.Time    `json:"valid_to"`
	Status       TariffStatus `json:"status,omitempty"`
	PublishedAt  *time.Time   `json:"published_at"`
	RetiredAt    *time.Time   `json:"retired_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	RoundingMode string       `json:"rounding_mode,omitempty"`
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
		ID:           tariff.ID,
		Code:         tariff.ShortCode,
		Description:  tariff.Description,
		Currency:     tariff.Currency,
		ValidFrom:    tariff.ValidFrom,
		ValidTo:      tariff.ValidTo,
		Status:       tariff.Status,
		PublishedAt:  tariff.PublishedAt,
		RetiredAt:    tariff.RetiredAt,
		CreatedAt:    tariff.CreatedAt,
		UpdatedAt:    tariff.UpdatedAt,
		RoundingMode: tariff.RoundingMode,
	}
}

type CreateWeightBreakRequest struct {
	TariffID         int64        `json:"tariff_id"`
	FromKg           float64      `json:"from_kg,omitempty"`
	ToKg             float64      `json:"to_kg,omitempty"`
	Zone             string       `json:"zone,omitempty"`
	RatePerKg        money.Amount `json:"rate_per_kg"`
	LoadingMeterRate money.Amount `json:"loading_meter_rate"`
	MinimumAmount    money.Amount `json:"minimum_amount"`
}

type WeightBreakResponse struct {
	ID               int64        `json:"id"`
	TariffID         int64        `json:"tariff_id"`
	FromKg           float64      `json:"from_kg,omitempty"`
	ToKg             float64      `json:"to_kg,omitempty"`
	Zone             string       `json:"zone,omitempty"`
	RatePerKg        money.Amount `json:"rate_per_kg"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	LoadingMeterRate money.Amount `json:"loading_meter_rate"`
	MinimumAmount    money.Amount `json:"minimum_amount"`
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
		ID:               weightBreak.ID,
		TariffID:         weightBreak.TariffID,
		FromKg:           weightBreak.FromKg,
		ToKg:             weightBreak.ToKg,
		Zone:             weightBreak.Zone,
		RatePerKg:        weightBreak.RatePerKg,
		CreatedAt:        weightBreak.CreatedAt,
		UpdatedAt:        weightBreak.UpdatedAt,
		LoadingMeterRate: weightBreak.LoadingMeterRate,
		MinimumAmount:    weightBreak.MinimumAmount,
	}
}

type CreateFuelSurchargeRequest struct {
	TariffID    int64        `json:"tariff_id"`
	ValidFrom   time.Time    `json:"valid_from"`
	Percent     float64      `json:"percent,omitempty"`
	DieselPrice money.Amount `json:"diesel_price"`
	Floor       float64      `json:"floor,omitempty"`
	Source      string       `json:"source,omitempty"`
}

type FuelSurchargeResponse struct {
	ID          int64        `json:"id"`
	TariffID    int64        `json:"tariff_id"`
	ValidFrom   time.Time    `json:"valid_from"`
	Percent     float64      `json:"percent,omitempty"`
	DieselPrice money.Amount `json:"diesel_price"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Floor       float64      `json:"floor,omitempty"`
	Source      string       `json:"source,omitempty"`
}

func (r *CreateFuelSurchargeRequest) Validate() error {
	if r.TariffID == 0 {
		return errors.New("tariff id is required")
	}
	if r.Percent < 0 {
		return errors.New("percent must not be negative")
	}
	return nil
}

func toFuelSurchargeResponse(fuelSurcharge *FuelSurcharge) FuelSurchargeResponse {
	return FuelSurchargeResponse{
		ID:          fuelSurcharge.ID,
		TariffID:    fuelSurcharge.TariffID,
		ValidFrom:   fuelSurcharge.ValidFrom,
		Percent:     fuelSurcharge.Percent,
		DieselPrice: fuelSurcharge.DieselPrice,
		CreatedAt:   fuelSurcharge.CreatedAt,
		UpdatedAt:   fuelSurcharge.UpdatedAt,
		Floor:       fuelSurcharge.Floor,
		Source:      fuelSurcharge.Source,
	}
}
