package shipment

import (
	"context"
)

// Store persists shipment data.
type Store interface {
	FindShipment(ctx context.Context, id int64) (*Shipment, error)
	ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error)
	ListShipmentsByStatus(ctx context.Context, status ShipmentStatus) ([]*Shipment, error)
	InsertShipment(ctx context.Context, shipment *Shipment) error
	UpdateShipment(ctx context.Context, shipment *Shipment) error
	FindTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error)
	ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error)
	ListTrackingEventsForShipment(ctx context.Context, shipmentID int64) ([]*TrackingEvent, error)
	InsertTrackingEvent(ctx context.Context, event *TrackingEvent) error
	UpdateTrackingEvent(ctx context.Context, event *TrackingEvent) error
}
