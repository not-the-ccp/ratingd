package shipment

import (
	"sort"
	"time"
)

// TrackingEvent holds the persisted state of a tracking event.
type TrackingEvent struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Location string `json:"location"`
	// ShipmentID references the owning shipment.
	ShipmentID int64     `json:"shipment_id"`
	OccurredAt time.Time `json:"occurred_at"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func sortTrackingEventsByOccurredAt(items []*TrackingEvent) {
	sort.Slice(items, func(i int, j int) bool {
		return items[i].OccurredAt.Before(items[j].OccurredAt)
	})
}

func (t *TrackingEvent) IsOccurredPassed(now time.Time) bool {
	return t.OccurredAt.Before(now)
}
