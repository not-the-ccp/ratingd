package tariff

import (
	"errors"
	"time"
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
