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

func (s *SQLStore) InsertFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO fuel_surcharges (tariff_id, valid_from, percent, diesel_price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, surcharge.TariffID, surcharge.ValidFrom, surcharge.Percent, surcharge.DieselPrice, surcharge.CreatedAt, surcharge.UpdatedAt).Scan(&surcharge.ID); err != nil {
		return err
	}
	return nil
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

func (s *SQLStore) FindWeightBreak(ctx context.Context, id int64) (*WeightBreak, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, loading_meter_rate, rate_per_kg, to_kg, from_kg, zone, tariff_id, created_at, updated_at FROM weight_breaks WHERE id = $1`, id)
	var breakEntry WeightBreak
	if err := row.Scan(&breakEntry.ID, &breakEntry.LoadingMeterRate, &breakEntry.RatePerKg, &breakEntry.ToKg, &breakEntry.FromKg, &breakEntry.Zone, &breakEntry.TariffID, &breakEntry.CreatedAt, &breakEntry.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query weight break %d: %w", id, err)
	}
	return &breakEntry, nil
}

func (s *SQLStore) ListWeightBreaksForTariff(ctx context.Context, tariffID int64) ([]*WeightBreak, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, loading_meter_rate, rate_per_kg, to_kg, from_kg, zone, tariff_id, created_at, updated_at FROM weight_breaks WHERE tariff_id = $1 ORDER BY id`, tariffID)
	if err != nil {
		return nil, fmt.Errorf("list weight breaks: %w", err)
	}
	defer rows.Close()

	var breakEntries []*WeightBreak
	for rows.Next() {
		var weightBreak WeightBreak
		if err := rows.Scan(&weightBreak.ID, &weightBreak.LoadingMeterRate, &weightBreak.RatePerKg, &weightBreak.ToKg, &weightBreak.FromKg, &weightBreak.Zone, &weightBreak.TariffID, &weightBreak.CreatedAt, &weightBreak.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan weight break: %w", err)
		}
		breakEntries = append(breakEntries, &weightBreak)
	}
	return breakEntries, rows.Err()
}

func (s *SQLStore) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, loading_meter_rate, rate_per_kg, to_kg, from_kg, zone, tariff_id, created_at, updated_at FROM weight_breaks ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list weight breaks: %w", err)
	}

	defer rows.Close()

	var breakEntries []*WeightBreak
	for rows.Next() {
		var breakEntry WeightBreak
		if err := rows.Scan(&breakEntry.ID, &breakEntry.LoadingMeterRate, &breakEntry.RatePerKg, &breakEntry.ToKg, &breakEntry.FromKg, &breakEntry.Zone, &breakEntry.TariffID, &breakEntry.CreatedAt, &breakEntry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan weight break: %w", err)
		}
		breakEntries = append(breakEntries, &breakEntry)
	}
	return breakEntries, rows.Err()
}

func (s *SQLStore) InsertWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO weight_breaks (loading_meter_rate, rate_per_kg, to_kg, from_kg, zone, tariff_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, breakVal.LoadingMeterRate, breakVal.RatePerKg, breakVal.ToKg, breakVal.FromKg, breakVal.Zone, breakVal.TariffID, breakVal.CreatedAt, breakVal.UpdatedAt).Scan(&breakVal.ID); err != nil {
		return err
	}
	return nil
}

func (s *SQLStore) UpdateWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	res, err := s.db.ExecContext(ctx, `UPDATE weight_breaks SET loading_meter_rate = $2, rate_per_kg = $3, to_kg = $4, from_kg = $5, zone = $6, tariff_id = $7, updated_at = $8 WHERE id = $1`, breakVal.ID, breakVal.LoadingMeterRate, breakVal.RatePerKg, breakVal.ToKg, breakVal.FromKg, breakVal.Zone, breakVal.TariffID, breakVal.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update weight break: %w", err)
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

func (s *SQLStore) FindFuelSurcharge(ctx context.Context, id int64) (*FuelSurcharge, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, tariff_id, valid_from, percent, diesel_price, created_at, updated_at FROM fuel_surcharges WHERE id = $1`, id)
	var surcharge FuelSurcharge
	if err := row.Scan(&surcharge.ID, &surcharge.TariffID, &surcharge.ValidFrom, &surcharge.Percent, &surcharge.DieselPrice, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query fuel surcharge %d: %w", id, err)
	}
	return &surcharge, nil
}

func (s *SQLStore) ListFuelSurcharges(ctx context.Context, limit int, offset int) ([]*FuelSurcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, tariff_id, valid_from, percent, diesel_price, created_at, updated_at FROM fuel_surcharges ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fuelSurcharges []*FuelSurcharge
	for rows.Next() {
		var surcharge FuelSurcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.TariffID, &surcharge.ValidFrom, &surcharge.Percent, &surcharge.DieselPrice, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan fuel surcharge: %w", err)
		}
		fuelSurcharges = append(fuelSurcharges, &surcharge)
	}
	return fuelSurcharges, rows.Err()
}

func (s *SQLStore) ListFuelSurchargesForTariff(ctx context.Context, tariffID int64) ([]*FuelSurcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, tariff_id, valid_from, percent, diesel_price, created_at, updated_at FROM fuel_surcharges WHERE tariff_id = $1 ORDER BY id`, tariffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surcharges []*FuelSurcharge
	for rows.Next() {
		var surcharge FuelSurcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.TariffID, &surcharge.ValidFrom, &surcharge.Percent, &surcharge.DieselPrice, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan fuel surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
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

func (s *SQLStore) UpdateFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	res, err := s.db.ExecContext(ctx, `UPDATE fuel_surcharges SET tariff_id = $2, valid_from = $3, percent = $4, diesel_price = $5, updated_at = $6 WHERE id = $1`, surcharge.ID, surcharge.TariffID, surcharge.ValidFrom, surcharge.Percent, surcharge.DieselPrice, surcharge.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update fuel surcharge: %w", err)
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
