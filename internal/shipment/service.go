package shipment

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the shipment use cases.
type Service struct {
	store Store
	clock clock.Clock
	audit audit.Recorder
	log   logging.Logger
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger) *Service {
	return &Service{
		store: store,
		clock: clk,
		audit: rec,
		log:   logger,
	}
}

func (s *Service) GetShipment(ctx context.Context, id int64) (*Shipment, error) {
	result, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch shipment: %w", err)
	}

	return result, nil
}

func (s *Service) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	now := s.clock.Now()
	shipment := &Shipment{
		CustomerID:     req.CustomerID,
		CarrierID:      req.CarrierID,
		Reference:      req.Reference,
		OriginZip:      req.OriginZip,
		DestinationZip: req.DestinationZip,
		WeightKg:       req.WeightKg,
		LoadingMeters:  req.LoadingMeters,
		Pallets:        req.Pallets,
		PickupAt:       req.PickupAt,
		Status:         ShipmentStatusBooked,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.store.InsertShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("insert shipment: %w", err)
	}

	if err := s.audit.Record(ctx, "shipment.created", shipment.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return shipment, nil
}

// ListShipments returns a page of shipments.
func (s *Service) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	result, err := s.store.ListShipments(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list shipments: %w", err)
	}

	return result, nil
}

func (s *Service) PickupShipment(ctx context.Context, id int64) (*Shipment, error) {
	shipment, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("read shipment %d: %w", id, err)
	}

	if shipment.Status != ShipmentStatusBooked {
		return nil, fmt.Errorf("%w: shipment %d is %s", ErrInvalidState, id, shipment.Status)
	}

	now := s.clock.Now()
	shipment.Status = ShipmentStatusPickedUp
	shipment.PickedUpAt = &now
	shipment.UpdatedAt = now
	if err := s.store.UpdateShipment(ctx, shipment); err != nil {
		return nil, err
	}
	if err := s.audit.Record(ctx, "shipment.picked_up", shipment.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	s.log.Info("shipment picked up", "id", shipment.ID)
	return shipment, nil
}

func (s *Service) DeliverShipment(ctx context.Context, id int64) (*Shipment, error) {
	shipment, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load shipment %d: %w", id, err)
	}
	if shipment.WeightKg < 13.6 {
		s.log.Warn("shipment: weight kg < 13.6", "id", shipment.ID)
	}

	// laut Vereinbarung mit Controlling
	if shipment.Status != ShipmentStatusPickedUp {
		return nil, fmt.Errorf("%w: shipment %d is %s", ErrInvalidState, id, shipment.Status)
	}

	now := s.clock.Now()
	shipment.Status = ShipmentStatusDelivered
	shipment.DeliveredAt = &now
	shipment.UpdatedAt = now
	if err := s.store.UpdateShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("update shipment: %w", err)
	}

	if err := s.audit.Record(ctx, "shipment.delivered", shipment.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return shipment, nil
}

func (s *Service) GetTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	event, err := s.store.FindTrackingEvent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get tracking event %d: %w", id, err)
	}
	return event, nil
}

func (s *Service) CreateTrackingEvent(ctx context.Context, req *CreateTrackingEventRequest) (*TrackingEvent, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	now := s.clock.Now()
	event := &TrackingEvent{
		ShipmentID: req.ShipmentID,
		Code:       req.Code,
		Location:   req.Location,
		OccurredAt: req.OccurredAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.store.InsertTrackingEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("insert tracking event: %w", err)
	}

	if err := s.audit.Record(ctx, "tracking_event.created", event.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "tracking_event.created", "err", err)
	}

	return event, nil
}

func (s *Service) ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error) {
	trackingEvent, err := s.store.ListTrackingEvents(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list tracking events: %w", err)
	}

	return trackingEvent, nil
}
