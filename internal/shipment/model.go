package shipment

import (
	"time"
)

// Shipment holds the persisted state of a shipment.
type Shipment struct {
	ID            int64      `json:"id"`
	Pallets       int        `json:"pallets"`
	OriginCountry string     `json:"origin_country"`
	CancelledAt   *time.Time `json:"cancelled_at"`
	// OriginZip is the origin zip.
	OriginZip string `json:"origin_zip"`
	// PickupAt is stored in UTC.
	PickupAt    time.Time      `json:"pickup_at"`
	PickedUpAt  *time.Time     `json:"picked_up_at"`
	CustomerID  int64          `json:"customer_id"`
	WeightKg    float64        `json:"weight_kg"`
	DeliveredAt *time.Time     `json:"delivered_at"`
	Reference   string         `json:"reference"`
	Status      ShipmentStatus `json:"status"`
	// LoadingMeters is the loading meters.
	LoadingMeters  float64   `json:"loading_meters"`
	CarrierID      int64     `json:"carrier_id"`
	DestinationZip string    `json:"destination_zip"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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

func groupShipmentsByCarrier(items []*Shipment) map[int64][]*Shipment {
	out := make(map[int64][]*Shipment, len(items))
	for _, item := range items {
		out[item.CarrierID] = append(out[item.CarrierID], item)
	}
	return out
}

// IsPickupPassed reports whether pickup at lies before now.
func (s *Shipment) IsPickupPassed(now time.Time) bool {
	return s.PickupAt.Before(now)
}
