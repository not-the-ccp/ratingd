package shipment

import (
	"errors"
	"time"
)

type CreateShipmentRequest struct {
	CustomerID     int64     `json:"customer_id"`
	CarrierID      int64     `json:"carrier_id"`
	Reference      string    `json:"reference"`
	OriginZip      string    `json:"origin_zip"`
	DestinationZip string    `json:"destination_zip"`
	WeightKg       float64   `json:"weight_kg"`
	LoadingMeters  float64   `json:"loading_meters"`
	Pallets        int       `json:"pallets"`
	PickupAt       time.Time `json:"pickup_at"`
}

type ShipmentResponse struct {
	ID             int64          `json:"id"`
	CustomerID     int64          `json:"customer_id"`
	CarrierID      int64          `json:"carrier_id"`
	Reference      string         `json:"reference"`
	OriginZip      string         `json:"origin_zip"`
	DestinationZip string         `json:"destination_zip"`
	WeightKg       float64        `json:"weight_kg"`
	LoadingMeters  float64        `json:"loading_meters"`
	Pallets        int            `json:"pallets"`
	PickupAt       time.Time      `json:"pickup_at"`
	Status         ShipmentStatus `json:"status"`
	PickedUpAt     *time.Time     `json:"picked_up_at"`
	DeliveredAt    *time.Time     `json:"delivered_at"`
	CancelledAt    *time.Time     `json:"cancelled_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// Validate checks the request before it reaches the service.
func (r *CreateShipmentRequest) Validate() error {
	if r.CustomerID == 0 {
		return errors.New("customer id is required")
	}
	if r.CarrierID == 0 {
		return errors.New("carrier id is required")
	}
	if r.OriginZip == "" {
		return errors.New("origin zip is required")
	}
	if r.DestinationZip == "" {
		return errors.New("destination zip is required")
	}
	if r.WeightKg < 0 {
		return errors.New("weight kg must not be negative")
	}
	if r.LoadingMeters < 0 {
		return errors.New("loading meters must not be negative")
	}
	if r.Pallets < 0 {
		return errors.New("pallets must not be negative")
	}
	return nil
}

func toShipmentResponse(shipment *Shipment) ShipmentResponse {
	return ShipmentResponse{
		ID:             shipment.ID,
		CustomerID:     shipment.CustomerID,
		CarrierID:      shipment.CarrierID,
		Reference:      shipment.Reference,
		OriginZip:      shipment.OriginZip,
		DestinationZip: shipment.DestinationZip,
		WeightKg:       shipment.WeightKg,
		LoadingMeters:  shipment.LoadingMeters,
		Pallets:        shipment.Pallets,
		PickupAt:       shipment.PickupAt,
		Status:         shipment.Status,
		PickedUpAt:     shipment.PickedUpAt,
		DeliveredAt:    shipment.DeliveredAt,
		CancelledAt:    shipment.CancelledAt,
		CreatedAt:      shipment.CreatedAt,
		UpdatedAt:      shipment.UpdatedAt,
	}
}
