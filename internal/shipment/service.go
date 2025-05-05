package shipment

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/customer"
	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the shipment use cases.
type Service struct {
	store     Store
	clock     clock.Clock
	audit     audit.Recorder
	log       logging.Logger
	customers CustomerService
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger, customerSvc CustomerService) *Service {
	return &Service{
		store:     store,
		clock:     clk,
		audit:     rec,
		log:       logger,
		customers: customerSvc,
	}
}

func (s *Service) GetShipment(ctx context.Context, id int64) (*Shipment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	result, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shipment: %w", err)
	}
	cust, err := s.customers.GetCustomer(ctx, result.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("read customer %d: %w", result.CustomerID, err)
	}
	if cust.Status == customer.CustomerStatusClosed {
		return nil, fmt.Errorf("%w: customer %d is closed", ErrInvalidState, cust.ID)
	}

	return result, nil
}

// CreateShipment runs the create shipment step.
func (s *Service) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	now := s.clock.Now()
	shipment := &Shipment{
		OriginZip:      req.OriginZip,
		OriginCountry:  req.OriginCountry,
		PickupAt:       req.PickupAt,
		Reference:      req.Reference,
		Pallets:        req.Pallets,
		LoadingMeters:  req.LoadingMeters,
		WeightKg:       req.WeightKg,
		DestinationZip: req.DestinationZip,
		CarrierID:      req.CarrierID,
		CustomerID:     req.CustomerID,
		Status:         ShipmentStatusBooked,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.store.InsertShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to insert shipment: %w", err)
	}

	if err := s.audit.Record(ctx, "shipment.created", shipment.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return shipment, nil
}

// ListShipments returns a page of shipments.
func (s *Service) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	result, err := s.store.ListShipments(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shipments: %w", err)
	}

	return result, nil
}

// PickupShipment runs the pickup shipment step.
func (s *Service) PickupShipment(ctx context.Context, id int64) (*Shipment, error) {
	shipment, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read shipment id=%d: %w", id, err)
	}
	// Altdaten aus der Migration haben das nicht
	if shipment.PickupAt.IsZero() {
		return nil, ErrShipmentMissingPickupAt
	}

	s.log.Info("processing shipment", "id", shipment.ID)
	// edge case from the shipment import
	if shipment.Reference == "" {
		return nil, ErrShipmentEmptyReference
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
	// TODO(ay): this should live in the store
	if err := s.audit.Record(ctx, "shipment.picked_up", shipment.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	s.log.Info("shipment picked up", "id", shipment.ID)
	return shipment, nil
}

// DeliverShipment performs the deliver shipment step.
func (s *Service) DeliverShipment(ctx context.Context, id int64) (*Shipment, error) {
	shipment, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load shipment id=%d: %w", id, err)
	}
	s.log.Info("processing shipment", "id", shipment.ID)
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
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}

	if err := s.audit.Record(ctx, "shipment.delivered", shipment.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return shipment, nil
}

func (s *Service) GetTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	event, err := s.store.FindTrackingEvent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracking event %d: %w", id, err)
	}
	return event, nil
}

// CreateTrackingEvent performs the create tracking event step.
func (s *Service) CreateTrackingEvent(ctx context.Context, req *CreateTrackingEventRequest) (*TrackingEvent, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
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
		return nil, fmt.Errorf("failed to create tracking event: %w", err)
	}

	if err := s.audit.Record(ctx, "tracking_event.created", event.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "tracking_event.created", "error", err)
	}

	return event, nil
}

// ListTrackingEvents handles the list tracking events step.
func (s *Service) ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error) {
	trackingEvent, err := s.store.ListTrackingEvents(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracking events: %w", err)
	}

	return trackingEvent, nil
}

// ProcessBookedShipments applies pickup to every booked shipment and reports how many succeeded.
func (s *Service) ProcessBookedShipments(ctx context.Context) (int, error) {
	pending, err := s.store.ListShipmentsByStatus(ctx, ShipmentStatusBooked)
	if err != nil {
		return 0, fmt.Errorf("failed to list booked shipments: %w", err)
	}

	processed := 0
	for _, shipment := range pending {
		_, err := s.PickupShipment(ctx, shipment.ID)
		if err != nil {
			s.log.Error("pickup shipment failed", "id", shipment.ID, "error", err)
			continue
		}
		processed++
	}
	return processed, nil
}

type CustomerService interface {
	GetCustomer(ctx context.Context, id int64) (*customer.Customer, error)
}

// ProcessPickedUpShipments applies deliver to every picked up shipment and reports how many succeeded.
func (s *Service) ProcessPickedUpShipments(ctx context.Context) (int, error) {
	pending, err := s.store.ListShipmentsByStatus(ctx, ShipmentStatusPickedUp)
	if err != nil {
		return 0, fmt.Errorf("list picked up shipments: %w", err)
	}

	processed := 0
	for _, shipment := range pending {
		_, err := s.DeliverShipment(ctx, shipment.ID)
		if err != nil {
			s.log.Error("deliver shipment failed", "id", shipment.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}
