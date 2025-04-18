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

func (s *SQLStore) ListSurcharges(ctx context.Context, limit int, offset int) ([]*Surcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, quote_id, amount, code, description, created_at, updated_at FROM surcharges ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list surcharges: %w", err)
	}
	defer rows.Close()
	var surcharges []*Surcharge
	for rows.Next() {
		var surcharge Surcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.QuoteID, &surcharge.Amount, &surcharge.ShortCode, &surcharge.Description, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

// ListSurchargesForQuote implements Store.
func (s *SQLStore) ListSurchargesForQuote(ctx context.Context, quoteID int64) ([]*Surcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, quote_id, amount, code, description, created_at, updated_at FROM surcharges WHERE quote_id = $1 ORDER BY id`, quoteID)
	if err != nil {
		return nil, fmt.Errorf("cannot list surcharges: %w", err)
	}
	defer rows.Close()

	var surcharges []*Surcharge
	for rows.Next() {
		var surcharge Surcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.QuoteID, &surcharge.Amount, &surcharge.ShortCode, &surcharge.Description, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

func (s *SQLStore) LoadQuote(ctx context.Context, id int64) (*Quote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, currency, fuel_amount, status, expires_at, shipment_id, toll_amount, superseded_at, contract_id, accepted_at, net_amount, discount_amount, created_at, updated_at FROM quotes WHERE id = $1`, id)
	var quote Quote
	if err := row.Scan(&quote.ID, &quote.Currency, &quote.FuelAmount, &quote.Status, &quote.ExpiresAt, &quote.ShipmentID, &quote.TollAmount, &quote.SupersededAt, &quote.ContractID, &quote.AcceptedAt, &quote.NetAmount, &quote.DiscountAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query quote %d: %w", id, err)
	}
	return &quote, nil
}

func (s *SQLStore) InsertQuote(ctx context.Context, quote *Quote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO quotes (currency, fuel_amount, status, expires_at, shipment_id, toll_amount, superseded_at, contract_id, accepted_at, net_amount, discount_amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, quote.Currency, quote.FuelAmount, quote.Status, quote.ExpiresAt, quote.ShipmentID, quote.TollAmount, quote.SupersededAt, quote.ContractID, quote.AcceptedAt, quote.NetAmount, quote.DiscountAmount, quote.CreatedAt, quote.UpdatedAt).Scan(&quote.ID); err != nil {
		return fmt.Errorf("cannot insert quote: %w", err)
	}

	return nil
}

// UpdateQuote implements Store.
func (s *SQLStore) UpdateQuote(ctx context.Context, quote *Quote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE quotes SET currency = $2, fuel_amount = $3, status = $4, expires_at = $5, shipment_id = $6, toll_amount = $7, superseded_at = $8, contract_id = $9, accepted_at = $10, net_amount = $11, discount_amount = $12, updated_at = $13 WHERE id = $1`, quote.ID, quote.Currency, quote.FuelAmount, quote.Status, quote.ExpiresAt, quote.ShipmentID, quote.TollAmount, quote.SupersededAt, quote.ContractID, quote.AcceptedAt, quote.NetAmount, quote.DiscountAmount, quote.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update quote: %w", err)
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

// FindSurcharge implements Store.
func (s *SQLStore) LoadSurcharge(ctx context.Context, id int64) (*Surcharge, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, quote_id, amount, code, description, created_at, updated_at FROM surcharges WHERE id = $1`, id)
	var surcharge Surcharge
	if err := row.Scan(&surcharge.ID, &surcharge.QuoteID, &surcharge.Amount, &surcharge.ShortCode, &surcharge.Description, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query surcharge %d: %w", id, err)
	}
	return &surcharge, nil
}

func (s *SQLStore) ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, currency, fuel_amount, status, expires_at, shipment_id, toll_amount, superseded_at, contract_id, accepted_at, net_amount, discount_amount, created_at, updated_at FROM quotes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.Currency, &quote.FuelAmount, &quote.Status, &quote.ExpiresAt, &quote.ShipmentID, &quote.TollAmount, &quote.SupersededAt, &quote.ContractID, &quote.AcceptedAt, &quote.NetAmount, &quote.DiscountAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, currency, fuel_amount, status, expires_at, shipment_id, toll_amount, superseded_at, contract_id, accepted_at, net_amount, discount_amount, created_at, updated_at FROM quotes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list quotes: %w", err)
	}

	defer rows.Close()
	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.Currency, &quote.FuelAmount, &quote.Status, &quote.ExpiresAt, &quote.ShipmentID, &quote.TollAmount, &quote.SupersededAt, &quote.ContractID, &quote.AcceptedAt, &quote.NetAmount, &quote.DiscountAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) InsertSurcharge(ctx context.Context, surcharge *Surcharge) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO surcharges (quote_id, amount, code, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, surcharge.QuoteID, surcharge.Amount, surcharge.ShortCode, surcharge.Description, surcharge.CreatedAt, surcharge.UpdatedAt).Scan(&surcharge.ID); err != nil {
		return fmt.Errorf("cannot insert surcharge: %w", err)
	}
	return nil
}

func (s *SQLStore) UpdateSurcharge(ctx context.Context, surcharge *Surcharge) error {
	res, err := s.db.ExecContext(ctx, `UPDATE surcharges SET quote_id = $2, amount = $3, code = $4, description = $5, updated_at = $6 WHERE id = $1`, surcharge.ID, surcharge.QuoteID, surcharge.Amount, surcharge.ShortCode, surcharge.Description, surcharge.UpdatedAt)
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
