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

func (s *SQLStore) ListWeightBreaksForTariff(ctx context.Context, tariffID int64) ([]*WeightBreak, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, from_kg, zone, tariff_id, to_kg, rate_per_kg, loading_meter_rate, created_at, updated_at FROM weight_breaks WHERE tariff_id = $1 ORDER BY id`, tariffID)
	if err != nil {
		return nil, fmt.Errorf("cannot list weight breaks: %w", err)
	}
	defer rows.Close()

	var breakEntries []*WeightBreak
	for rows.Next() {
		var weightBreak WeightBreak
		if err := rows.Scan(&weightBreak.ID, &weightBreak.FromKg, &weightBreak.Zone, &weightBreak.TariffID, &weightBreak.ToKg, &weightBreak.RatePerKg, &weightBreak.LoadingMeterRate, &weightBreak.CreatedAt, &weightBreak.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan weight break: %w", err)
		}
		breakEntries = append(breakEntries, &weightBreak)
	}
	return breakEntries, rows.Err()
}

func (s *SQLStore) InsertFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO fuel_surcharges (floor, diesel_price, valid_from, tariff_id, percent, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, surcharge.Floor, surcharge.DieselPrice, surcharge.ValidFrom, surcharge.TariffID, surcharge.Percent, surcharge.CreatedAt, surcharge.UpdatedAt).Scan(&surcharge.ID); err != nil {
		return err
	}
	return nil
}

func (s *SQLStore) ListFuelSurchargesForTariff(ctx context.Context, tariffID int64) ([]*FuelSurcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, floor, diesel_price, valid_from, tariff_id, percent, created_at, updated_at FROM fuel_surcharges WHERE tariff_id = $1 ORDER BY id`, tariffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surcharges []*FuelSurcharge
	for rows.Next() {
		var surcharge FuelSurcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.Floor, &surcharge.DieselPrice, &surcharge.ValidFrom, &surcharge.TariffID, &surcharge.Percent, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan fuel surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

// ListTariffsByStatus implements Store.
func (s *SQLStore) ListTariffsByStatus(ctx context.Context, status TariffStatus) ([]*Tariff, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, status, rounding_mode, code, valid_from, currency, retired_at, description, published_at, valid_to, created_at, updated_at FROM tariffs WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list tariffs: %w", err)
	}
	defer rows.Close()
	var tariffs []*Tariff
	for rows.Next() {
		var tariff Tariff
		if err := rows.Scan(&tariff.ID, &tariff.Status, &tariff.RoundingMode, &tariff.Code, &tariff.ValidFrom, &tariff.Currency, &tariff.RetiredAt, &tariff.Description, &tariff.PublishedAt, &tariff.ValidTo, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan tariff: %w", err)
		}
		tariffs = append(tariffs, &tariff)
	}
	return tariffs, rows.Err()
}

func (s *SQLStore) InsertTariff(ctx context.Context, tariff *Tariff) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO tariffs (status, rounding_mode, code, valid_from, currency, retired_at, description, published_at, valid_to, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, tariff.Status, tariff.RoundingMode, tariff.Code, tariff.ValidFrom, tariff.Currency, tariff.RetiredAt, tariff.Description, tariff.PublishedAt, tariff.ValidTo, tariff.CreatedAt, tariff.UpdatedAt).Scan(&tariff.ID); err != nil {
		return fmt.Errorf("cannot insert tariff: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateTariff(ctx context.Context, tariff *Tariff) error {
	res, err := s.db.ExecContext(ctx, `UPDATE tariffs SET status = $2, rounding_mode = $3, code = $4, valid_from = $5, currency = $6, retired_at = $7, description = $8, published_at = $9, valid_to = $10, updated_at = $11 WHERE id = $1`, tariff.ID, tariff.Status, tariff.RoundingMode, tariff.Code, tariff.ValidFrom, tariff.Currency, tariff.RetiredAt, tariff.Description, tariff.PublishedAt, tariff.ValidTo, tariff.UpdatedAt)
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
	row := s.db.QueryRowContext(ctx, `SELECT id, from_kg, zone, tariff_id, to_kg, rate_per_kg, loading_meter_rate, created_at, updated_at FROM weight_breaks WHERE id = $1`, id)
	var breakEntry WeightBreak
	if err := row.Scan(&breakEntry.ID, &breakEntry.FromKg, &breakEntry.Zone, &breakEntry.TariffID, &breakEntry.ToKg, &breakEntry.RatePerKg, &breakEntry.LoadingMeterRate, &breakEntry.CreatedAt, &breakEntry.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query weight break %d: %w", id, err)
	}
	return &breakEntry, nil
}

// FindTariff implements Store.
func (s *SQLStore) FindTariff(ctx context.Context, id int64) (*Tariff, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, status, rounding_mode, code, valid_from, currency, retired_at, description, published_at, valid_to, created_at, updated_at FROM tariffs WHERE id = $1`, id)
	var tariff Tariff
	if err := row.Scan(&tariff.ID, &tariff.Status, &tariff.RoundingMode, &tariff.Code, &tariff.ValidFrom, &tariff.Currency, &tariff.RetiredAt, &tariff.Description, &tariff.PublishedAt, &tariff.ValidTo, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query tariff %d: %w", id, err)
	}
	return &tariff, nil
}

func (s *SQLStore) ListFuelSurcharges(ctx context.Context, limit int, offset int) ([]*FuelSurcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, floor, diesel_price, valid_from, tariff_id, percent, created_at, updated_at FROM fuel_surcharges ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fuelSurcharges []*FuelSurcharge
	for rows.Next() {
		var surcharge FuelSurcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.Floor, &surcharge.DieselPrice, &surcharge.ValidFrom, &surcharge.TariffID, &surcharge.Percent, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan fuel surcharge: %w", err)
		}
		fuelSurcharges = append(fuelSurcharges, &surcharge)
	}
	return fuelSurcharges, rows.Err()
}

func (s *SQLStore) InsertWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO weight_breaks (from_kg, zone, tariff_id, to_kg, rate_per_kg, loading_meter_rate, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, breakVal.FromKg, breakVal.Zone, breakVal.TariffID, breakVal.ToKg, breakVal.RatePerKg, breakVal.LoadingMeterRate, breakVal.CreatedAt, breakVal.UpdatedAt).Scan(&breakVal.ID); err != nil {
		return err
	}
	return nil
}

func (s *SQLStore) UpdateWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	res, err := s.db.ExecContext(ctx, `UPDATE weight_breaks SET from_kg = $2, zone = $3, tariff_id = $4, to_kg = $5, rate_per_kg = $6, loading_meter_rate = $7, updated_at = $8 WHERE id = $1`, breakVal.ID, breakVal.FromKg, breakVal.Zone, breakVal.TariffID, breakVal.ToKg, breakVal.RatePerKg, breakVal.LoadingMeterRate, breakVal.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update weight break: %w", err)
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
	row := s.db.QueryRowContext(ctx, `SELECT id, floor, diesel_price, valid_from, tariff_id, percent, created_at, updated_at FROM fuel_surcharges WHERE id = $1`, id)
	var surcharge FuelSurcharge
	if err := row.Scan(&surcharge.ID, &surcharge.Floor, &surcharge.DieselPrice, &surcharge.ValidFrom, &surcharge.TariffID, &surcharge.Percent, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query fuel surcharge %d: %w", id, err)
	}
	return &surcharge, nil
}

func (s *SQLStore) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, from_kg, zone, tariff_id, to_kg, rate_per_kg, loading_meter_rate, created_at, updated_at FROM weight_breaks ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list weight breaks: %w", err)
	}

	defer rows.Close()

	var breakEntries []*WeightBreak
	for rows.Next() {
		var breakEntry WeightBreak
		if err := rows.Scan(&breakEntry.ID, &breakEntry.FromKg, &breakEntry.Zone, &breakEntry.TariffID, &breakEntry.ToKg, &breakEntry.RatePerKg, &breakEntry.LoadingMeterRate, &breakEntry.CreatedAt, &breakEntry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan weight break: %w", err)
		}
		breakEntries = append(breakEntries, &breakEntry)
	}
	return breakEntries, rows.Err()
}

// ListTariffs implements Store.
func (s *SQLStore) ListTariffs(ctx context.Context, limit int, offset int) ([]*Tariff, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, status, rounding_mode, code, valid_from, currency, retired_at, description, published_at, valid_to, created_at, updated_at FROM tariffs ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list tariffs: %w", err)
	}

	defer rows.Close()

	var tariffs []*Tariff
	for rows.Next() {
		var tariff Tariff
		if err := rows.Scan(&tariff.ID, &tariff.Status, &tariff.RoundingMode, &tariff.Code, &tariff.ValidFrom, &tariff.Currency, &tariff.RetiredAt, &tariff.Description, &tariff.PublishedAt, &tariff.ValidTo, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan tariff: %w", err)
		}
		tariffs = append(tariffs, &tariff)
	}
	return tariffs, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) UpdateFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	res, err := s.db.ExecContext(ctx, `UPDATE fuel_surcharges SET floor = $2, diesel_price = $3, valid_from = $4, tariff_id = $5, percent = $6, updated_at = $7 WHERE id = $1`, surcharge.ID, surcharge.Floor, surcharge.DieselPrice, surcharge.ValidFrom, surcharge.TariffID, surcharge.Percent, surcharge.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update fuel surcharge: %w", err)
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
