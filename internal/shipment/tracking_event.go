package shipment

import (
	"time"
)

// TrackingEvent holds the persisted state of a tracking event.
type TrackingEvent struct {
	ID       int64  `json:"id"`
	Source   string `json:"source"`
	Location string `json:"location"`
	Code     string `json:"code"`
	// ShipmentID references the owning shipment.
	ShipmentID int64     `json:"shipment_id"`
	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *TrackingEvent) IsOccurredPassed(now time.Time) bool {
	return t.OccurredAt.Before(now)
}
