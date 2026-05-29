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

func (s *SQLStore) ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, shipment_id, toll_amount, accepted_at, superseded_at, expires_at, net_amount, status, fuel_amount, currency, contract_id, discount_amount, created_at, updated_at FROM quotes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.ShipmentID, &quote.TollAmount, &quote.AcceptedAt, &quote.SupersededAt, &quote.ExpiresAt, &quote.NetAmount, &quote.Status, &quote.FuelAmount, &quote.Currency, &quote.ContractID, &quote.DiscountAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) InsertSurcharge(ctx context.Context, surcharge *Surcharge) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO surcharges (code, quote_id, description, amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, surcharge.ShortCode, surcharge.QuoteID, surcharge.Description, surcharge.Amount, surcharge.CreatedAt, surcharge.UpdatedAt).Scan(&surcharge.ID); err != nil {
		return fmt.Errorf("failed to insert surcharge: %w", err)
	}
	return nil
}

// UpdateQuote implements Store.
func (s *SQLStore) UpdateQuote(ctx context.Context, quote *Quote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE quotes SET shipment_id = $2, toll_amount = $3, accepted_at = $4, superseded_at = $5, expires_at = $6, net_amount = $7, status = $8, fuel_amount = $9, currency = $10, contract_id = $11, discount_amount = $12, updated_at = $13 WHERE id = $1`, quote.ID, quote.ShipmentID, quote.TollAmount, quote.AcceptedAt, quote.SupersededAt, quote.ExpiresAt, quote.NetAmount, quote.Status, quote.FuelAmount, quote.Currency, quote.ContractID, quote.DiscountAmount, quote.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update quote: %w", err)
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

func (s *SQLStore) LoadQuote(ctx context.Context, id int64) (*Quote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, shipment_id, toll_amount, accepted_at, superseded_at, expires_at, net_amount, status, fuel_amount, currency, contract_id, discount_amount, created_at, updated_at FROM quotes WHERE id = $1`, id)
	var quote Quote
	if err := row.Scan(&quote.ID, &quote.ShipmentID, &quote.TollAmount, &quote.AcceptedAt, &quote.SupersededAt, &quote.ExpiresAt, &quote.NetAmount, &quote.Status, &quote.FuelAmount, &quote.Currency, &quote.ContractID, &quote.DiscountAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query quote %d: %w", id, err)
	}
	return &quote, nil
}

// FindSurcharge implements Store.
func (s *SQLStore) LoadSurcharge(ctx context.Context, id int64) (*Surcharge, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, code, quote_id, description, amount, created_at, updated_at FROM surcharges WHERE id = $1`, id)
	var surcharge Surcharge
	if err := row.Scan(&surcharge.ID, &surcharge.ShortCode, &surcharge.QuoteID, &surcharge.Description, &surcharge.Amount, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query surcharge %d: %w", id, err)
	}
	return &surcharge, nil
}

func (s *SQLStore) InsertQuote(ctx context.Context, quote *Quote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO quotes (shipment_id, toll_amount, accepted_at, superseded_at, expires_at, net_amount, status, fuel_amount, currency, contract_id, discount_amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, quote.ShipmentID, quote.TollAmount, quote.AcceptedAt, quote.SupersededAt, quote.ExpiresAt, quote.NetAmount, quote.Status, quote.FuelAmount, quote.Currency, quote.ContractID, quote.DiscountAmount, quote.CreatedAt, quote.UpdatedAt).Scan(&quote.ID); err != nil {
		return fmt.Errorf("failed to insert quote: %w", err)
	}

	return nil
}

func (s *SQLStore) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, shipment_id, toll_amount, accepted_at, superseded_at, expires_at, net_amount, status, fuel_amount, currency, contract_id, discount_amount, created_at, updated_at FROM quotes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}

	defer rows.Close()
	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.ShipmentID, &quote.TollAmount, &quote.AcceptedAt, &quote.SupersededAt, &quote.ExpiresAt, &quote.NetAmount, &quote.Status, &quote.FuelAmount, &quote.Currency, &quote.ContractID, &quote.DiscountAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) ListSurcharges(ctx context.Context, limit int, offset int) ([]*Surcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, quote_id, description, amount, created_at, updated_at FROM surcharges ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list surcharges: %w", err)
	}
	defer rows.Close()
	var surcharges []*Surcharge
	for rows.Next() {
		var surcharge Surcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.ShortCode, &surcharge.QuoteID, &surcharge.Description, &surcharge.Amount, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

func (s *SQLStore) UpdateSurcharge(ctx context.Context, surcharge *Surcharge) error {
	res, err := s.db.ExecContext(ctx, `UPDATE surcharges SET code = $2, quote_id = $3, description = $4, amount = $5, updated_at = $6 WHERE id = $1`, surcharge.ID, surcharge.ShortCode, surcharge.QuoteID, surcharge.Description, surcharge.Amount, surcharge.UpdatedAt)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

// ListSurchargesForQuote implements Store.
func (s *SQLStore) ListSurchargesForQuote(ctx context.Context, quoteID int64) ([]*Surcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, quote_id, description, amount, created_at, updated_at FROM surcharges WHERE quote_id = $1 ORDER BY id`, quoteID)
	if err != nil {
		return nil, fmt.Errorf("failed to list surcharges: %w", err)
	}
	defer rows.Close()

	var surcharges []*Surcharge
	for rows.Next() {
		var surcharge Surcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.ShortCode, &surcharge.QuoteID, &surcharge.Description, &surcharge.Amount, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}
