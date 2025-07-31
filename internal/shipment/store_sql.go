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

func (s *SQLStore) FindTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, code, location, shipment_id, occurred_at, created_at, updated_at FROM tracking_events WHERE id = $1`, id)
	var event TrackingEvent
	if err := row.Scan(&event.ID, &event.Code, &event.Location, &event.ShipmentID, &event.OccurredAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query tracking event %d: %w", id, err)
	}
	return &event, nil
}

func (s *SQLStore) ListTrackingEventsForShipment(ctx context.Context, shipmentID int64) ([]*TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, location, shipment_id, occurred_at, created_at, updated_at FROM tracking_events WHERE shipment_id = $1 ORDER BY id`, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking events: %w", err)
	}

	defer rows.Close()
	var events []*TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.ID, &event.Code, &event.Location, &event.ShipmentID, &event.OccurredAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tracking event: %w", err)
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLStore) ListShipmentsByStatus(ctx context.Context, status ShipmentStatus) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, destination_zip, reference, status, pallets, origin_country, loading_meters, carrier_id, cancelled_at, delivered_at, customer_id, pickup_at, picked_up_at, weight_kg, origin_zip, created_at, updated_at FROM shipments WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.DestinationZip, &shipment.Reference, &shipment.Status, &shipment.Pallets, &shipment.OriginCountry, &shipment.LoadingMeters, &shipment.CarrierID, &shipment.CancelledAt, &shipment.DeliveredAt, &shipment.CustomerID, &shipment.PickupAt, &shipment.PickedUpAt, &shipment.WeightKg, &shipment.OriginZip, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, location, shipment_id, occurred_at, created_at, updated_at FROM tracking_events ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking events: %w", err)
	}

	defer rows.Close()

	var events []*TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.ID, &event.Code, &event.Location, &event.ShipmentID, &event.OccurredAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tracking event: %w", err)
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLStore) UpdateShipment(ctx context.Context, shipment *Shipment) error {
	res, err := s.db.ExecContext(ctx, `UPDATE shipments SET destination_zip = $2, reference = $3, status = $4, pallets = $5, origin_country = $6, loading_meters = $7, carrier_id = $8, cancelled_at = $9, delivered_at = $10, customer_id = $11, pickup_at = $12, picked_up_at = $13, weight_kg = $14, origin_zip = $15, updated_at = $16 WHERE id = $1`, shipment.ID, shipment.DestinationZip, shipment.Reference, shipment.Status, shipment.Pallets, shipment.OriginCountry, shipment.LoadingMeters, shipment.CarrierID, shipment.CancelledAt, shipment.DeliveredAt, shipment.CustomerID, shipment.PickupAt, shipment.PickedUpAt, shipment.WeightKg, shipment.OriginZip, shipment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, destination_zip, reference, status, pallets, origin_country, loading_meters, carrier_id, cancelled_at, delivered_at, customer_id, pickup_at, picked_up_at, weight_kg, origin_zip, created_at, updated_at FROM shipments ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shipments: %w", err)
	}
	defer rows.Close()

	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.DestinationZip, &shipment.Reference, &shipment.Status, &shipment.Pallets, &shipment.OriginCountry, &shipment.LoadingMeters, &shipment.CarrierID, &shipment.CancelledAt, &shipment.DeliveredAt, &shipment.CustomerID, &shipment.PickupAt, &shipment.PickedUpAt, &shipment.WeightKg, &shipment.OriginZip, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) InsertTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO tracking_events (code, location, shipment_id, occurred_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, event.Code, event.Location, event.ShipmentID, event.OccurredAt, event.CreatedAt, event.UpdatedAt).Scan(&event.ID); err != nil {
		return fmt.Errorf("failed to insert tracking event: %w", err)
	}
	return nil
}

func (s *SQLStore) InsertShipment(ctx context.Context, shipment *Shipment) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO shipments (destination_zip, reference, status, pallets, origin_country, loading_meters, carrier_id, cancelled_at, delivered_at, customer_id, pickup_at, picked_up_at, weight_kg, origin_zip, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`, shipment.DestinationZip, shipment.Reference, shipment.Status, shipment.Pallets, shipment.OriginCountry, shipment.LoadingMeters, shipment.CarrierID, shipment.CancelledAt, shipment.DeliveredAt, shipment.CustomerID, shipment.PickupAt, shipment.PickedUpAt, shipment.WeightKg, shipment.OriginZip, shipment.CreatedAt, shipment.UpdatedAt).Scan(&shipment.ID); err != nil {
		return fmt.Errorf("failed to insert shipment: %w", err)
	}
	return nil
}

// UpdateTrackingEvent implements Store.
func (s *SQLStore) UpdateTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	res, err := s.db.ExecContext(ctx, `UPDATE tracking_events SET code = $2, location = $3, shipment_id = $4, occurred_at = $5, updated_at = $6 WHERE id = $1`, event.ID, event.Code, event.Location, event.ShipmentID, event.OccurredAt, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update tracking event: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) FindShipment(ctx context.Context, id int64) (*Shipment, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, destination_zip, reference, status, pallets, origin_country, loading_meters, carrier_id, cancelled_at, delivered_at, customer_id, pickup_at, picked_up_at, weight_kg, origin_zip, created_at, updated_at FROM shipments WHERE id = $1`, id)
	var shipment Shipment
	if err := row.Scan(&shipment.ID, &shipment.DestinationZip, &shipment.Reference, &shipment.Status, &shipment.Pallets, &shipment.OriginCountry, &shipment.LoadingMeters, &shipment.CarrierID, &shipment.CancelledAt, &shipment.DeliveredAt, &shipment.CustomerID, &shipment.PickupAt, &shipment.PickedUpAt, &shipment.WeightKg, &shipment.OriginZip, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query shipment %d: %w", id, err)
	}
	return &shipment, nil
}
