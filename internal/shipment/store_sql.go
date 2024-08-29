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
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, created_at, updated_at FROM shipments WHERE id = $1`, id)
	var shipment Shipment
	if err := row.Scan(&shipment.ID, &shipment.CustomerID, &shipment.CarrierID, &shipment.Reference, &shipment.OriginZip, &shipment.DestinationZip, &shipment.WeightKg, &shipment.LoadingMeters, &shipment.Pallets, &shipment.PickupAt, &shipment.Status, &shipment.PickedUpAt, &shipment.DeliveredAt, &shipment.CancelledAt, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query shipment %d: %w", id, err)
	}
	return &shipment, nil
}

func (s *SQLStore) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, created_at, updated_at FROM shipments ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list shipments: %w", err)
	}
	defer rows.Close()

	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.CustomerID, &shipment.CarrierID, &shipment.Reference, &shipment.OriginZip, &shipment.DestinationZip, &shipment.WeightKg, &shipment.LoadingMeters, &shipment.Pallets, &shipment.PickupAt, &shipment.Status, &shipment.PickedUpAt, &shipment.DeliveredAt, &shipment.CancelledAt, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) ListShipmentsByStatus(ctx context.Context, status ShipmentStatus) ([]*Shipment, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, created_at, updated_at FROM shipments WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shipments []*Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.CustomerID, &shipment.CarrierID, &shipment.Reference, &shipment.OriginZip, &shipment.DestinationZip, &shipment.WeightKg, &shipment.LoadingMeters, &shipment.Pallets, &shipment.PickupAt, &shipment.Status, &shipment.PickedUpAt, &shipment.DeliveredAt, &shipment.CancelledAt, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shipment: %w", err)
		}
		shipments = append(shipments, &shipment)
	}
	return shipments, rows.Err()
}

func (s *SQLStore) InsertShipment(ctx context.Context, shipment *Shipment) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO shipments (customer_id, carrier_id, reference, origin_zip, destination_zip, weight_kg, loading_meters, pallets, pickup_at, status, picked_up_at, delivered_at, cancelled_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`, shipment.CustomerID, shipment.CarrierID, shipment.Reference, shipment.OriginZip, shipment.DestinationZip, shipment.WeightKg, shipment.LoadingMeters, shipment.Pallets, shipment.PickupAt, shipment.Status, shipment.PickedUpAt, shipment.DeliveredAt, shipment.CancelledAt, shipment.CreatedAt, shipment.UpdatedAt).Scan(&shipment.ID); err != nil {
		return fmt.Errorf("insert shipment: %w", err)
	}
	return nil
}

func (s *SQLStore) UpdateShipment(ctx context.Context, shipment *Shipment) error {
	res, err := s.db.ExecContext(ctx, `UPDATE shipments SET customer_id = $2, carrier_id = $3, reference = $4, origin_zip = $5, destination_zip = $6, weight_kg = $7, loading_meters = $8, pallets = $9, pickup_at = $10, status = $11, picked_up_at = $12, delivered_at = $13, cancelled_at = $14, updated_at = $15 WHERE id = $1`, shipment.ID, shipment.CustomerID, shipment.CarrierID, shipment.Reference, shipment.OriginZip, shipment.DestinationZip, shipment.WeightKg, shipment.LoadingMeters, shipment.Pallets, shipment.PickupAt, shipment.Status, shipment.PickedUpAt, shipment.DeliveredAt, shipment.CancelledAt, shipment.UpdatedAt)
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
