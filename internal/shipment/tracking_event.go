package shipment

import (
	"time"
)

// TrackingEvent holds the persisted state of a tracking event.
type TrackingEvent struct {
	ID         int64     `json:"id"`
	OccurredAt time.Time `json:"occurred_at"`
	Location   string    `json:"location"`
	Code       string    `json:"code"`
	Source     string    `json:"source"`
	// ShipmentID references the owning shipment.
	ShipmentID int64     `json:"shipment_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *TrackingEvent) IsOccurredPassed(now time.Time) bool {
	return t.OccurredAt.Before(now)
}
