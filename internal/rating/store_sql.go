package rating

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

func (s *SQLStore) GetQuote(ctx context.Context, id int64) (*Quote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, shipment_id, contract_id, net_amount, fuel_amount, currency, expires_at, status, accepted_at, superseded_at, created_at, updated_at FROM quotes WHERE id = $1`, id)
	var quote Quote
	if err := row.Scan(&quote.ID, &quote.ShipmentID, &quote.ContractID, &quote.NetAmount, &quote.FuelAmount, &quote.Currency, &quote.ExpiresAt, &quote.Status, &quote.AcceptedAt, &quote.SupersededAt, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query quote %d: %w", id, err)
	}
	return &quote, nil
}

func (s *SQLStore) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, shipment_id, contract_id, net_amount, fuel_amount, currency, expires_at, status, accepted_at, superseded_at, created_at, updated_at FROM quotes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list quotes: %w", err)
	}

	defer rows.Close()
	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.ShipmentID, &quote.ContractID, &quote.NetAmount, &quote.FuelAmount, &quote.Currency, &quote.ExpiresAt, &quote.Status, &quote.AcceptedAt, &quote.SupersededAt, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, shipment_id, contract_id, net_amount, fuel_amount, currency, expires_at, status, accepted_at, superseded_at, created_at, updated_at FROM quotes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.ShipmentID, &quote.ContractID, &quote.NetAmount, &quote.FuelAmount, &quote.Currency, &quote.ExpiresAt, &quote.Status, &quote.AcceptedAt, &quote.SupersededAt, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) InsertQuote(ctx context.Context, quote *Quote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO quotes (shipment_id, contract_id, net_amount, fuel_amount, currency, expires_at, status, accepted_at, superseded_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, quote.ShipmentID, quote.ContractID, quote.NetAmount, quote.FuelAmount, quote.Currency, quote.ExpiresAt, quote.Status, quote.AcceptedAt, quote.SupersededAt, quote.CreatedAt, quote.UpdatedAt).Scan(&quote.ID); err != nil {
		return fmt.Errorf("insert quote: %w", err)
	}

	return nil
}

// UpdateQuote implements Store.
func (s *SQLStore) UpdateQuote(ctx context.Context, quote *Quote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE quotes SET shipment_id = $2, contract_id = $3, net_amount = $4, fuel_amount = $5, currency = $6, expires_at = $7, status = $8, accepted_at = $9, superseded_at = $10, updated_at = $11 WHERE id = $1`, quote.ID, quote.ShipmentID, quote.ContractID, quote.NetAmount, quote.FuelAmount, quote.Currency, quote.ExpiresAt, quote.Status, quote.AcceptedAt, quote.SupersededAt, quote.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update quote: %w", err)
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
