package shipment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SQLStore implements Store on top of PostgreSQL.
type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

// UpdateTrackingEvent implements Store.
func (s *SQLStore) UpdateTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	res, err := s.db.ExecContext(ctx, `UPDATE tracking_events SET source = $2, location = $3, code = $4, shipment_id = $5, occurred_at = $6, updated_at = $7 WHERE id = $1`, event.ID, event.Source, event.Location, event.Code, event.ShipmentID, event.OccurredAt, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update tracking event: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) InsertTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO tracking_events (source, location, code, shipment_id, occurred_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, event.Source, event.Location, event.Code, event.ShipmentID, event.OccurredAt, event.CreatedAt, event.UpdatedAt).Scan(&event.ID); err != nil {
		return fmt.Errorf("cannot insert tracking event: %w", err)
	}
	return nil
}

func (s *SQLStore) ListShipmentsByStatus(ctx context.Context, status ShipmentStatus) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, loading_meters, hazard_class, status, weight_kg, pickup_at, carrier_id, reference, destination_zip, pallets, delivered_at, origin_country, customer_id, cancelled_at, origin_zip, picked_up_at, created_at, updated_at FROM shipments WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.LoadingMeters, &shipment.HazardClass, &shipment.Status, &shipment.WeightKg, &shipment.PickupAt, &shipment.CarrierID, &shipment.Reference, &shipment.DestinationZip, &shipment.Pallets, &shipment.DeliveredAt, &shipment.OriginCountry, &shipment.CustomerID, &shipment.CancelledAt, &shipment.OriginZip, &shipment.PickedUpAt, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) UpdateShipment(ctx context.Context, shipment *Shipment) error {
	res, err := s.db.ExecContext(ctx, `UPDATE shipments SET loading_meters = $2, hazard_class = $3, status = $4, weight_kg = $5, pickup_at = $6, carrier_id = $7, reference = $8, destination_zip = $9, pallets = $10, delivered_at = $11, origin_country = $12, customer_id = $13, cancelled_at = $14, origin_zip = $15, picked_up_at = $16, updated_at = $17 WHERE id = $1`, shipment.ID, shipment.LoadingMeters, shipment.HazardClass, shipment.Status, shipment.WeightKg, shipment.PickupAt, shipment.CarrierID, shipment.Reference, shipment.DestinationZip, shipment.Pallets, shipment.DeliveredAt, shipment.OriginCountry, shipment.CustomerID, shipment.CancelledAt, shipment.OriginZip, shipment.PickedUpAt, shipment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update shipment: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, loading_meters, hazard_class, status, weight_kg, pickup_at, carrier_id, reference, destination_zip, pallets, delivered_at, origin_country, customer_id, cancelled_at, origin_zip, picked_up_at, created_at, updated_at FROM shipments ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list shipments: %w", err)
	}
	defer rows.Close()

	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.LoadingMeters, &shipment.HazardClass, &shipment.Status, &shipment.WeightKg, &shipment.PickupAt, &shipment.CarrierID, &shipment.Reference, &shipment.DestinationZip, &shipment.Pallets, &shipment.DeliveredAt, &shipment.OriginCountry, &shipment.CustomerID, &shipment.CancelledAt, &shipment.OriginZip, &shipment.PickedUpAt, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) FetchShipment(ctx context.Context, id int64) (*Shipment, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, loading_meters, hazard_class, status, weight_kg, pickup_at, carrier_id, reference, destination_zip, pallets, delivered_at, origin_country, customer_id, cancelled_at, origin_zip, picked_up_at, created_at, updated_at FROM shipments WHERE id = $1`, id)
	var shipment Shipment
	if err := row.Scan(&shipment.ID, &shipment.LoadingMeters, &shipment.HazardClass, &shipment.Status, &shipment.WeightKg, &shipment.PickupAt, &shipment.CarrierID, &shipment.Reference, &shipment.DestinationZip, &shipment.Pallets, &shipment.DeliveredAt, &shipment.OriginCountry, &shipment.CustomerID, &shipment.CancelledAt, &shipment.OriginZip, &shipment.PickedUpAt, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query shipment %d: %w", id, err)
	}
	return &shipment, nil
}

func (s *SQLStore) ListTrackingEventsForShipment(ctx context.Context, shipmentID int64) ([]*TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, source, location, code, shipment_id, occurred_at, created_at, updated_at FROM tracking_events WHERE shipment_id = $1 ORDER BY id`, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("cannot list tracking events: %w", err)
	}

	defer rows.Close()
	var events []*TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.ID, &event.Source, &event.Location, &event.Code, &event.ShipmentID, &event.OccurredAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan tracking event: %w", err)
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLStore) ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, source, location, code, shipment_id, occurred_at, created_at, updated_at FROM tracking_events ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list tracking events: %w", err)
	}

	defer rows.Close()

	var events []*TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.ID, &event.Source, &event.Location, &event.Code, &event.ShipmentID, &event.OccurredAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan tracking event: %w", err)
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLStore) InsertShipment(ctx context.Context, shipment *Shipment) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO shipments (loading_meters, hazard_class, status, weight_kg, pickup_at, carrier_id, reference, destination_zip, pallets, delivered_at, origin_country, customer_id, cancelled_at, origin_zip, picked_up_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`, shipment.LoadingMeters, shipment.HazardClass, shipment.Status, shipment.WeightKg, shipment.PickupAt, shipment.CarrierID, shipment.Reference, shipment.DestinationZip, shipment.Pallets, shipment.DeliveredAt, shipment.OriginCountry, shipment.CustomerID, shipment.CancelledAt, shipment.OriginZip, shipment.PickedUpAt, shipment.CreatedAt, shipment.UpdatedAt).Scan(&shipment.ID); err != nil {
		return fmt.Errorf("cannot insert shipment: %w", err)
	}
	return nil
}

func (s *SQLStore) FetchTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, source, location, code, shipment_id, occurred_at, created_at, updated_at FROM tracking_events WHERE id = $1`, id)
	var event TrackingEvent
	if err := row.Scan(&event.ID, &event.Source, &event.Location, &event.Code, &event.ShipmentID, &event.OccurredAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query tracking event %d: %w", id, err)
	}
	return &event, nil
}
