package carrier

import (
	"errors"
	"time"
)

type CreateCarrierRequest struct {
	Name                string    `json:"name"`
	Code                string    `json:"code"`
	Country             string    `json:"country"`
	VATNumber           string    `json:"vat_number"`
	PaymentTermDays     int       `json:"payment_term_days"`
	InsuranceValidUntil time.Time `json:"insurance_valid_until"`
	FleetSize           int       `json:"fleet_size"`
	IBAN                string    `json:"iban"`
}

type CarrierResponse struct {
	ID                  int64         `json:"id"`
	Name                string        `json:"name"`
	Code                string        `json:"code"`
	Country             string        `json:"country"`
	VATNumber           string        `json:"vat_number"`
	PaymentTermDays     int           `json:"payment_term_days"`
	Status              CarrierStatus `json:"status"`
	ApprovedAt          *time.Time    `json:"approved_at"`
	SuspendedAt         *time.Time    `json:"suspended_at"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	InsuranceValidUntil time.Time     `json:"insurance_valid_until"`
	FleetSize           int           `json:"fleet_size"`
	IBAN                string        `json:"iban"`
}

// Validate checks the request before it reaches the service.
func (r *CreateCarrierRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Country == "" {
		return errors.New("country is required")
	}
	if r.VATNumber == "" {
		return errors.New("vat number is required")
	}
	if r.PaymentTermDays < 0 {
		return errors.New("payment term days must not be negative")
	}
	return nil
}

func toCarrierResponse(carrier *Carrier) CarrierResponse {
	return CarrierResponse{
		ID:                  carrier.ID,
		Name:                carrier.DisplayName,
		Code:                carrier.ShortCode,
		Country:             carrier.Country,
		VATNumber:           carrier.VATNumber,
		PaymentTermDays:     carrier.PaymentTermDays,
		Status:              carrier.Status,
		ApprovedAt:          carrier.ApprovedAt,
		SuspendedAt:         carrier.SuspendedAt,
		CreatedAt:           carrier.CreatedAt,
		UpdatedAt:           carrier.UpdatedAt,
		InsuranceValidUntil: carrier.InsuranceValidUntil,
		FleetSize:           carrier.FleetSize,
		IBAN:                carrier.IBAN,
	}
}
