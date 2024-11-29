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

func (s *SQLStore) ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, currency, net_amount, expires_at, status, accepted_at, contract_id, fuel_amount, superseded_at, shipment_id, toll_amount, created_at, updated_at FROM quotes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.Currency, &quote.NetAmount, &quote.ExpiresAt, &quote.Status, &quote.AcceptedAt, &quote.ContractID, &quote.FuelAmount, &quote.SupersededAt, &quote.ShipmentID, &quote.TollAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, currency, net_amount, expires_at, status, accepted_at, contract_id, fuel_amount, superseded_at, shipment_id, toll_amount, created_at, updated_at FROM quotes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list quotes: %w", err)
	}

	defer rows.Close()
	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID, &quote.Currency, &quote.NetAmount, &quote.ExpiresAt, &quote.Status, &quote.AcceptedAt, &quote.ContractID, &quote.FuelAmount, &quote.SupersededAt, &quote.ShipmentID, &quote.TollAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes, rows.Err()
}

func (s *SQLStore) GetQuote(ctx context.Context, id int64) (*Quote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, currency, net_amount, expires_at, status, accepted_at, contract_id, fuel_amount, superseded_at, shipment_id, toll_amount, created_at, updated_at FROM quotes WHERE id = $1`, id)
	var quote Quote
	if err := row.Scan(&quote.ID, &quote.Currency, &quote.NetAmount, &quote.ExpiresAt, &quote.Status, &quote.AcceptedAt, &quote.ContractID, &quote.FuelAmount, &quote.SupersededAt, &quote.ShipmentID, &quote.TollAmount, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query quote %d: %w", id, err)
	}
	return &quote, nil
}

func (s *SQLStore) InsertQuote(ctx context.Context, quote *Quote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO quotes (currency, net_amount, expires_at, status, accepted_at, contract_id, fuel_amount, superseded_at, shipment_id, toll_amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`, quote.Currency, quote.NetAmount, quote.ExpiresAt, quote.Status, quote.AcceptedAt, quote.ContractID, quote.FuelAmount, quote.SupersededAt, quote.ShipmentID, quote.TollAmount, quote.CreatedAt, quote.UpdatedAt).Scan(&quote.ID); err != nil {
		return fmt.Errorf("insert quote: %w", err)
	}

	return nil
}

// UpdateQuote implements Store.
func (s *SQLStore) UpdateQuote(ctx context.Context, quote *Quote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE quotes SET currency = $2, net_amount = $3, expires_at = $4, status = $5, accepted_at = $6, contract_id = $7, fuel_amount = $8, superseded_at = $9, shipment_id = $10, toll_amount = $11, updated_at = $12 WHERE id = $1`, quote.ID, quote.Currency, quote.NetAmount, quote.ExpiresAt, quote.Status, quote.AcceptedAt, quote.ContractID, quote.FuelAmount, quote.SupersededAt, quote.ShipmentID, quote.TollAmount, quote.UpdatedAt)
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

// FindSurcharge implements Store.
func (s *SQLStore) FindSurcharge(ctx context.Context, id int64) (*Surcharge, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, quote_id, amount, code, description, created_at, updated_at FROM surcharges WHERE id = $1`, id)
	var surcharge Surcharge
	if err := row.Scan(&surcharge.ID, &surcharge.QuoteID, &surcharge.Amount, &surcharge.Code, &surcharge.Description, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query surcharge %d: %w", id, err)
	}
	return &surcharge, nil
}

func (s *SQLStore) ListSurcharges(ctx context.Context, limit int, offset int) ([]*Surcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, quote_id, amount, code, description, created_at, updated_at FROM surcharges ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list surcharges: %w", err)
	}
	defer rows.Close()
	var surcharges []*Surcharge
	for rows.Next() {
		var surcharge Surcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.QuoteID, &surcharge.Amount, &surcharge.Code, &surcharge.Description, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

// ListSurchargesForQuote implements Store.
func (s *SQLStore) ListSurchargesForQuote(ctx context.Context, quoteID int64) ([]*Surcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, quote_id, amount, code, description, created_at, updated_at FROM surcharges WHERE quote_id = $1 ORDER BY id`, quoteID)
	if err != nil {
		return nil, fmt.Errorf("list surcharges: %w", err)
	}
	defer rows.Close()

	var surcharges []*Surcharge
	for rows.Next() {
		var surcharge Surcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.QuoteID, &surcharge.Amount, &surcharge.Code, &surcharge.Description, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

func (s *SQLStore) InsertSurcharge(ctx context.Context, surcharge *Surcharge) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO surcharges (quote_id, amount, code, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, surcharge.QuoteID, surcharge.Amount, surcharge.Code, surcharge.Description, surcharge.CreatedAt, surcharge.UpdatedAt).Scan(&surcharge.ID); err != nil {
		return fmt.Errorf("insert surcharge: %w", err)
	}
	return nil
}

func (s *SQLStore) UpdateSurcharge(ctx context.Context, surcharge *Surcharge) error {
	res, err := s.db.ExecContext(ctx, `UPDATE surcharges SET quote_id = $2, amount = $3, code = $4, description = $5, updated_at = $6 WHERE id = $1`, surcharge.ID, surcharge.QuoteID, surcharge.Amount, surcharge.Code, surcharge.Description, surcharge.UpdatedAt)
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
