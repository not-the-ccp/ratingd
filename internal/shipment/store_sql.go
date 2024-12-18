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

func (s *SQLStore) FindShipment(ctx context.Context, id int64) (*Shipment, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, origin_country, created_at, updated_at FROM shipments WHERE id = $1`, id)
	var shipment Shipment
	if err := row.Scan(&shipment.ID, &shipment.CustomerID, &shipment.CarrierID, &shipment.Reference, &shipment.OriginZip, &shipment.DestinationZip, &shipment.WeightKg, &shipment.LoadingMeters, &shipment.Pallets, &shipment.PickupAt, &shipment.Status, &shipment.PickedUpAt, &shipment.DeliveredAt, &shipment.CancelledAt, &shipment.OriginCountry, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query shipment %d: %w", id, err)
	}
	return &shipment, nil
}

func (s *SQLStore) FindTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, location, occurred_at, shipment_id, code, created_at, updated_at FROM tracking_events WHERE id = $1`, id)
	var event TrackingEvent
	if err := row.Scan(&event.ID, &event.Location, &event.OccurredAt, &event.ShipmentID, &event.Code, &event.CreatedAt, &event.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query tracking event %d: %w", id, err)
	}
	return &event, nil
}

func (s *SQLStore) ListShipmentsByStatus(ctx context.Context, status ShipmentStatus) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, origin_country, created_at, updated_at FROM shipments WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.CustomerID, &shipment.CarrierID, &shipment.Reference, &shipment.OriginZip, &shipment.DestinationZip, &shipment.WeightKg, &shipment.LoadingMeters, &shipment.Pallets, &shipment.PickupAt, &shipment.Status, &shipment.PickedUpAt, &shipment.DeliveredAt, &shipment.CancelledAt, &shipment.OriginCountry, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) InsertShipment(ctx context.Context, shipment *Shipment) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO shipments (customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, origin_country, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`, shipment.CustomerID, shipment.CarrierID, shipment.Reference, shipment.OriginZip, shipment.DestinationZip, shipment.WeightKg, shipment.LoadingMeters, shipment.Pallets, shipment.PickupAt, shipment.Status, shipment.PickedUpAt, shipment.DeliveredAt, shipment.CancelledAt, shipment.OriginCountry, shipment.CreatedAt, shipment.UpdatedAt).Scan(&shipment.ID); err != nil {
		return fmt.Errorf("insert shipment: %w", err)
	}
	return nil
}

func (s *SQLStore) UpdateShipment(ctx context.Context, shipment *Shipment) error {
	res, err := s.db.ExecContext(ctx, `UPDATE shipments SET customer_id = $2, carrier_id = $3, reference = $4, origin_zip = $5, destination_zip = $6, weight_kg = $7, loading_meters = $8, pallets = $9, pickup_at = $10, status = $11, picked_up_at = $12, delivered_at = $13, cancelled_at = $14, origin_country = $15, updated_at = $16 WHERE id = $1`, shipment.ID, shipment.CustomerID, shipment.CarrierID, shipment.Reference, shipment.OriginZip, shipment.DestinationZip, shipment.WeightKg, shipment.LoadingMeters, shipment.Pallets, shipment.PickupAt, shipment.Status, shipment.PickedUpAt, shipment.DeliveredAt, shipment.CancelledAt, shipment.OriginCountry, shipment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update shipment: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) InsertTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO tracking_events (location, occurred_at, shipment_id, code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, event.Location, event.OccurredAt, event.ShipmentID, event.Code, event.CreatedAt, event.UpdatedAt).Scan(&event.ID); err != nil {
		return fmt.Errorf("insert tracking event: %w", err)
	}
	return nil
}

func (s *SQLStore) ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, location, occurred_at, shipment_id, code, created_at, updated_at FROM tracking_events ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list tracking events: %w", err)
	}

	defer rows.Close()

	var events []*TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.ID, &event.Location, &event.OccurredAt, &event.ShipmentID, &event.Code, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan tracking event: %w", err)
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLStore) ListTrackingEventsForShipment(ctx context.Context, shipmentID int64) ([]*TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, location, occurred_at, shipment_id, code, created_at, updated_at FROM tracking_events WHERE shipment_id = $1 ORDER BY id`, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("list tracking events: %w", err)
	}

	defer rows.Close()
	var events []*TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.ID, &event.Location, &event.OccurredAt, &event.ShipmentID, &event.Code, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan tracking event: %w", err)
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLStore) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, origin_country, created_at, updated_at FROM shipments ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list shipments: %w", err)
	}
	defer rows.Close()

	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.CustomerID, &shipment.CarrierID, &shipment.Reference, &shipment.OriginZip, &shipment.DestinationZip, &shipment.WeightKg, &shipment.LoadingMeters, &shipment.Pallets, &shipment.PickupAt, &shipment.Status, &shipment.PickedUpAt, &shipment.DeliveredAt, &shipment.CancelledAt, &shipment.OriginCountry, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

// UpdateTrackingEvent implements Store.
func (s *SQLStore) UpdateTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	res, err := s.db.ExecContext(ctx, `UPDATE tracking_events SET location = $2, occurred_at = $3, shipment_id = $4, code = $5, updated_at = $6 WHERE id = $1`, event.ID, event.Location, event.OccurredAt, event.ShipmentID, event.Code, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update tracking event: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
