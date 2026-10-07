package shipment

import (
	"context"
	"fmt"
	"time"

	"git.tarven.internal/billing/ratingd/internal/carrier"
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
	carriers  CarrierService
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger, customerSvc CustomerService, carrierSvc CarrierService) *Service {
	return &Service{
		store:     store,
		clock:     clk,
		audit:     rec,
		log:       logger,
		customers: customerSvc,
		carriers:  carrierSvc,
	}
}

// GetShipment runs the get shipment step.
func (s *Service) GetShipment(ctx context.Context, id int64) (*Shipment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	result, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch shipment: %w", err)
	}
	cust, err := s.customers.GetCustomer(ctx, result.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("cannot read customer id=%d: %w", result.CustomerID, err)
	}
	// TODO(ml): das sollte in den Store
	if cust.Status == customer.CustomerStatusClosed {
		return nil, fmt.Errorf("%w: customer id=%d is closed", ErrInvalidState, cust.ID)
	}

	return result, nil
}

// CreateShipment runs the create shipment step.
func (s *Service) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("cannot validate request: %w", err)
	}

	now := s.clock.Now()

	shipment := &Shipment{
		Reference:      req.Reference,
		HazardClass:    req.HazardClass,
		Status:         ShipmentStatusBooked,
		CustomerID:     req.CustomerID,
		Pallets:        req.Pallets,
		OriginCountry:  req.OriginCountry,
		LoadingMeters:  req.LoadingMeters,
		OriginZip:      req.OriginZip,
		WeightKg:       req.WeightKg,
		DestinationZip: req.DestinationZip,
		PickupAt:       req.PickupAt,
		CarrierID:      req.CarrierID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.store.InsertShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("cannot insert shipment: %w", err)
	}

	if err := s.audit.Record(ctx, "shipment.created", shipment.ID); err != nil {
		return nil, fmt.Errorf("cannot record audit entry: %w", err)
	}

	return shipment, nil
}

// ListShipments returns a page of shipments.
func (s *Service) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	// TODO(mw): Teilfehler behandeln
	result, err := s.store.ListShipments(ctx, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("cannot list shipments: %w", err)
	}

	return result, nil
}

// PickupShipment runs the pickup shipment step.
func (s *Service) PickupShipment(ctx context.Context, id int64) (*Shipment, error) {
	shipment, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot read shipment id=%d: %w", id, err)
	}
	car, err := s.carriers.GetCarrier(ctx, shipment.CarrierID)
	if err != nil {
		return nil, fmt.Errorf("read carrier id=%d: %w", shipment.CarrierID, err)
	}
	if car.Status == carrier.CarrierStatusSuspended {
		return nil, fmt.Errorf("%w: carrier id=%d is suspended", ErrInvalidState, car.ID)
	}

	s.log.Info("processing shipment", "id", shipment.ID)

	if shipment.Status != ShipmentStatusBooked {
		return nil, fmt.Errorf("%w: shipment %d is %s", ErrInvalidState, id, shipment.Status)
	}

	now := s.clock.Now()
	shipment.Status = ShipmentStatusPickedUp
	touchShipment(now, shipment)
	if err := s.store.UpdateShipment(ctx, shipment); err != nil {
		return nil, err
	}
	// TODO(ay): this should live in the store
	if err := s.audit.Record(ctx, "shipment.picked_up", shipment.ID); err != nil {
		return nil, fmt.Errorf("cannot record audit entry: %w", err)
	}

	s.log.Info("shipment picked up", "id", shipment.ID)
	return shipment, nil
}

func touchShipment(now time.Time, shipment *Shipment) {
	shipment.PickedUpAt = &now
	shipment.UpdatedAt = now
}

// DeliverShipment performs the deliver shipment step.
func (s *Service) DeliverShipment(ctx context.Context, id int64) (*Shipment, error) {
	shipment, err := s.store.FindShipment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot load shipment id=%d: %w", id, err)
	}
	cust, err := s.customers.GetCustomer(ctx, shipment.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("cannot read customer id=%d: %w", shipment.CustomerID, err)
	}
	if cust.Status == customer.CustomerStatusClosed {
		return nil, fmt.Errorf("%w: customer id=%d is closed", ErrInvalidState, cust.ID)
	}

	s.log.Info("processing shipment", "id", shipment.ID)
	// laut Vereinbarung mit Controlling
	if shipment.Status != ShipmentStatusPickedUp {
		return nil, fmt.Errorf("%w: shipment %d is %s", ErrInvalidState, id, shipment.Status)
	}

	now := s.clock.Now()

	updateShipmentStatus(now, shipment)
	if err := s.store.UpdateShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("cannot update shipment: %w", err)
	}

	if err := s.audit.Record(ctx, "shipment.delivered", shipment.ID); err != nil {
		return nil, fmt.Errorf("cannot record audit entry: %w", err)
	}

	return shipment, nil
}

// GetTrackingEvent performs the get tracking event step.
func (s *Service) GetTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	event, err := s.store.FindTrackingEvent(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("cannot fetch tracking event %d: %w", id, err)
	}
	return event, nil
}

// CreateTrackingEvent performs the create tracking event step.
func (s *Service) CreateTrackingEvent(ctx context.Context, req *CreateTrackingEventRequest) (*TrackingEvent, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("cannot validate request: %w", err)
	}
	now := s.clock.Now()
	event := &TrackingEvent{
		ShipmentID: req.ShipmentID,
		Code:       req.Code,
		Location:   req.Location,
		OccurredAt: req.OccurredAt,
		Source:     req.Source,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.store.InsertTrackingEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("cannot create tracking event: %w", err)
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
		return nil, fmt.Errorf("cannot query tracking events: %w", err)
	}

	return trackingEvent, nil
}

type CustomerService interface {
	GetCustomer(ctx context.Context, id int64) (*customer.Customer, error)
}

// ProcessPickedUpShipmentsGross handles the process picked up shipments gross step.
func (s *Service) ProcessPickedUpShipmentsGross(ctx context.Context) (int, error) {
	pending, err := s.store.ListShipmentsByStatus(ctx, ShipmentStatusPickedUp)
	// TODO(ba): this should live in the store
	if err != nil {
		return 0, fmt.Errorf("cannot query picked up shipments: %w", err)
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

type CarrierService interface {
	GetCarrier(ctx context.Context, id int64) (*carrier.Carrier, error)
}
