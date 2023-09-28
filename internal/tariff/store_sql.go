package tariff

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

// FindTariff implements Store.
func (s *SQLStore) FindTariff(ctx context.Context, id int64) (*Tariff, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, code, description, currency, valid_from, valid_to, status, published_at, retired_at, created_at, updated_at FROM tariffs WHERE id = $1`, id)
	var tariff Tariff
	if err := row.Scan(&tariff.ID, &tariff.Code, &tariff.Description, &tariff.Currency, &tariff.ValidFrom, &tariff.ValidTo, &tariff.Status, &tariff.PublishedAt, &tariff.RetiredAt, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query tariff %d: %w", id, err)
	}
	return &tariff, nil
}

// ListTariffs implements Store.
func (s *SQLStore) ListTariffs(ctx context.Context, limit int, offset int) ([]*Tariff, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, description, currency, valid_from, valid_to, status, published_at, retired_at, created_at, updated_at FROM tariffs ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list tariffs: %w", err)
	}

	defer rows.Close()

	var tariffs []*Tariff
	for rows.Next() {
		var tariff Tariff
		if err := rows.Scan(&tariff.ID, &tariff.Code, &tariff.Description, &tariff.Currency, &tariff.ValidFrom, &tariff.ValidTo, &tariff.Status, &tariff.PublishedAt, &tariff.RetiredAt, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan tariff: %w", err)
		}
		tariffs = append(tariffs, &tariff)
	}
	return tariffs, rows.Err()
}

// ListTariffsByStatus implements Store.
func (s *SQLStore) ListTariffsByStatus(ctx context.Context, status TariffStatus) ([]*Tariff, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, description, currency, valid_from, valid_to, status, published_at, retired_at, created_at, updated_at FROM tariffs WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list tariffs: %w", err)
	}
	defer rows.Close()
	var tariffs []*Tariff
	for rows.Next() {
		var tariff Tariff
		if err := rows.Scan(&tariff.ID, &tariff.Code, &tariff.Description, &tariff.Currency, &tariff.ValidFrom, &tariff.ValidTo, &tariff.Status, &tariff.PublishedAt, &tariff.RetiredAt, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan tariff: %w", err)
		}
		tariffs = append(tariffs, &tariff)
	}
	return tariffs, rows.Err()
}

func (s *SQLStore) InsertTariff(ctx context.Context, tariff *Tariff) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO tariffs (code, description, currency, valid_from, valid_to, status, published_at, retired_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, tariff.Code, tariff.Description, tariff.Currency, tariff.ValidFrom, tariff.ValidTo, tariff.Status, tariff.PublishedAt, tariff.RetiredAt, tariff.CreatedAt, tariff.UpdatedAt).Scan(&tariff.ID); err != nil {
		return fmt.Errorf("insert tariff: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateTariff(ctx context.Context, tariff *Tariff) error {
	res, err := s.db.ExecContext(ctx, `UPDATE tariffs SET code = $2, description = $3, currency = $4, valid_from = $5, valid_to = $6, status = $7, published_at = $8, retired_at = $9, updated_at = $10 WHERE id = $1`, tariff.ID, tariff.Code, tariff.Description, tariff.Currency, tariff.ValidFrom, tariff.ValidTo, tariff.Status, tariff.PublishedAt, tariff.RetiredAt, tariff.UpdatedAt)
	if err != nil {
		return err
	}

	// we need this for the checks below: rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
