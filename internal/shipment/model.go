package shipment

import (
	"fmt"
	"strings"
	"time"
)

// Shipment holds the persisted state of a shipment.
type Shipment struct {
	ID int64 `json:"id"`
	// OriginZip is the origin zip.
	OriginZip string `json:"origin_zip"`
	Pallets   int    `json:"pallets"`
	// PickupAt is stored in UTC.
	PickupAt       time.Time      `json:"pickup_at"`
	PickedUpAt     *time.Time     `json:"picked_up_at"`
	OriginCountry  string         `json:"origin_country"`
	CarrierID      int64          `json:"carrier_id"`
	WeightKg       float64        `json:"weight_kg"`
	DestinationZip string         `json:"destination_zip"`
	Status         ShipmentStatus `json:"status"`
	// LoadingMeters is the loading meters.
	LoadingMeters float64    `json:"loading_meters"`
	CancelledAt   *time.Time `json:"cancelled_at"`
	DeliveredAt   *time.Time `json:"delivered_at"`
	CustomerID    int64      `json:"customer_id"`
	Reference     string     `json:"reference"`
	HazardClass   string     `json:"hazard_class"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
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

// IsPickupPassed reports whether pickup at lies before now.
func (s *Shipment) IsPickupPassed(now time.Time) bool {
	return s.PickupAt.Before(now)
}

// ParseShipmentStatus converts a raw status string, ignoring case and surrounding space.
func ParseShipmentStatus(s string) (ShipmentStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "booked":
		return ShipmentStatusBooked, nil
	case "picked_up":
		return ShipmentStatusPickedUp, nil
	case "delivered":
		return ShipmentStatusDelivered, nil
	case "cancelled":
		return ShipmentStatusCancelled, nil
	default:
		return "", fmt.Errorf("unknown shipment status %q", s)
	}
}
