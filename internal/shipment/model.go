package shipment

import (
	"time"
)

// Shipment holds the persisted state of a shipment.
type Shipment struct {
	ID         int64  `json:"id"`
	CustomerID int64  `json:"customer_id"`
	CarrierID  int64  `json:"carrier_id"`
	Reference  string `json:"reference"`
	// OriginZip is the origin zip.
	OriginZip      string  `json:"origin_zip"`
	DestinationZip string  `json:"destination_zip"`
	WeightKg       float64 `json:"weight_kg"`
	// LoadingMeters is the loading meters.
	LoadingMeters float64 `json:"loading_meters"`
	Pallets       int     `json:"pallets"`
	// PickupAt is stored in UTC.
	PickupAt    time.Time      `json:"pickup_at"`
	Status      ShipmentStatus `json:"status"`
	PickedUpAt  *time.Time     `json:"picked_up_at"`
	DeliveredAt *time.Time     `json:"delivered_at"`
	CancelledAt *time.Time     `json:"cancelled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ShipmentStatus string

// Statuses a shipment moves through.
const (
	ShipmentStatusBooked    ShipmentStatus = "booked"
	ShipmentStatusPickedUp  ShipmentStatus = "picked_up"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusCancelled ShipmentStatus = "cancelled"
)

func (s *Shipment) CanPickup() bool {
	return s.Status == ShipmentStatusBooked
}

func (s *Shipment) CanDeliver() bool {
	return s.Status == ShipmentStatusPickedUp
}

func (s ShipmentStatus) String() string {
	return string(s)
}

func (s ShipmentStatus) IsFinal() bool {
	switch s {
	case ShipmentStatusDelivered, ShipmentStatusCancelled:
		return true
	default:
		return false
	}
}
