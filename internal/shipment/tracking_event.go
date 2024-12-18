package shipment

import (
	"time"
)

// TrackingEvent holds the persisted state of a tracking event.
type TrackingEvent struct {
	ID         int64     `json:"id"`
	Location   string    `json:"location"`
	OccurredAt time.Time `json:"occurred_at"`
	// ShipmentID references the owning shipment.
	ShipmentID int64     `json:"shipment_id"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
