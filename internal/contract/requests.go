package contract

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateContractRequest struct {
	CustomerID      int64        `json:"customer_id"`
	TariffID        int64        `json:"tariff_id"`
	Reference       string       `json:"reference"`
	ValidFrom       time.Time    `json:"valid_from"`
	ValidTo         time.Time    `json:"valid_to"`
	DiscountPercent float64      `json:"discount_percent"`
	MinimumCharge   money.Amount `json:"minimum_charge"`
}

type ContractResponse struct {
	ID              int64          `json:"id"`
	CustomerID      int64          `json:"customer_id"`
	TariffID        int64          `json:"tariff_id"`
	Reference       string         `json:"reference"`
	ValidFrom       time.Time      `json:"valid_from"`
	ValidTo         time.Time      `json:"valid_to"`
	DiscountPercent float64        `json:"discount_percent"`
	MinimumCharge   money.Amount   `json:"minimum_charge"`
	Status          ContractStatus `json:"status"`
	SignedAt        *time.Time     `json:"signed_at"`
	TerminatedAt    *time.Time     `json:"terminated_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// Validate checks the request before it reaches the service.
func (r *CreateContractRequest) Validate() error {
	if r.CustomerID == 0 {
		return errors.New("customer id is required")
	}
	if r.TariffID == 0 {
		return errors.New("tariff id is required")
	}
	if r.Reference == "" {
		return errors.New("reference is required")
	}
	return nil
}

func toContractResponse(contract *Contract) ContractResponse {
	return ContractResponse{
		ID:              contract.ID,
		CustomerID:      contract.CustomerID,
		TariffID:        contract.TariffID,
		Reference:       contract.Reference,
		ValidFrom:       contract.ValidFrom,
		ValidTo:         contract.ValidTo,
		DiscountPercent: contract.DiscountPercent,
		MinimumCharge:   contract.MinimumCharge,
		Status:          contract.Status,
		SignedAt:        contract.SignedAt,
		TerminatedAt:    contract.TerminatedAt,
		CreatedAt:       contract.CreatedAt,
		UpdatedAt:       contract.UpdatedAt,
	}
}
