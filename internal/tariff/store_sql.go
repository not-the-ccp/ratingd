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

func (s *SQLStore) UpdateTariff(ctx context.Context, tariff *Tariff) error {
	res, err := s.db.ExecContext(ctx, `UPDATE tariffs SET retired_at = $2, published_at = $3, currency = $4, description = $5, valid_to = $6, status = $7, code = $8, rounding_mode = $9, base_zone = $10, valid_from = $11, region = $12, minimum_charge = $13, updated_at = $14 WHERE id = $1`, tariff.ID, tariff.RetiredAt, tariff.PublishedAt, tariff.Currency, tariff.Description, tariff.ValidTo, tariff.Status, tariff.ShortCode, tariff.RoundingMode, tariff.BaseZone, tariff.ValidFrom, tariff.Region, tariff.MinimumCharge, tariff.UpdatedAt)
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

func (s *SQLStore) InsertFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO fuel_surcharges (valid_from, diesel_price, floor, tariff_id, source, percent, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, surcharge.ValidFrom, surcharge.DieselPrice, surcharge.Floor, surcharge.TariffID, surcharge.Source, surcharge.Percent, surcharge.CreatedAt, surcharge.UpdatedAt).Scan(&surcharge.ID); err != nil {
		return err
	}
	return nil
}

func (s *SQLStore) FetchWeightBreak(ctx context.Context, id int64) (*WeightBreak, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, minimum_amount, from_kg, zone, rate_per_kg, tariff_id, to_kg, loading_meter_rate, created_at, updated_at FROM weight_breaks WHERE id = $1`, id)
	var breakEntry WeightBreak
	if err := row.Scan(&breakEntry.ID, &breakEntry.MinimumAmount, &breakEntry.FromKg, &breakEntry.Zone, &breakEntry.RatePerKg, &breakEntry.TariffID, &breakEntry.ToKg, &breakEntry.LoadingMeterRate, &breakEntry.CreatedAt, &breakEntry.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query weight break %d: %w", id, err)
	}
	return &breakEntry, nil
}

func (s *SQLStore) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, minimum_amount, from_kg, zone, rate_per_kg, tariff_id, to_kg, loading_meter_rate, created_at, updated_at FROM weight_breaks ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list weight breaks: %w", err)
	}

	defer rows.Close()

	var breakEntries []*WeightBreak
	for rows.Next() {
		var breakEntry WeightBreak
		if err := rows.Scan(&breakEntry.ID, &breakEntry.MinimumAmount, &breakEntry.FromKg, &breakEntry.Zone, &breakEntry.RatePerKg, &breakEntry.TariffID, &breakEntry.ToKg, &breakEntry.LoadingMeterRate, &breakEntry.CreatedAt, &breakEntry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan weight break: %w", err)
		}
		breakEntries = append(breakEntries, &breakEntry)
	}
	return breakEntries, rows.Err()
}

func (s *SQLStore) UpdateFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	res, err := s.db.ExecContext(ctx, `UPDATE fuel_surcharges SET valid_from = $2, diesel_price = $3, floor = $4, tariff_id = $5, source = $6, percent = $7, updated_at = $8 WHERE id = $1`, surcharge.ID, surcharge.ValidFrom, surcharge.DieselPrice, surcharge.Floor, surcharge.TariffID, surcharge.Source, surcharge.Percent, surcharge.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update fuel surcharge: %w", err)
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

func (s *SQLStore) ListFuelSurchargesForTariff(ctx context.Context, tariffID int64) ([]*FuelSurcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, valid_from, diesel_price, floor, tariff_id, source, percent, created_at, updated_at FROM fuel_surcharges WHERE tariff_id = $1 ORDER BY id`, tariffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surcharges []*FuelSurcharge
	for rows.Next() {
		var surcharge FuelSurcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.ValidFrom, &surcharge.DieselPrice, &surcharge.Floor, &surcharge.TariffID, &surcharge.Source, &surcharge.Percent, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan fuel surcharge: %w", err)
		}
		surcharges = append(surcharges, &surcharge)
	}
	return surcharges, rows.Err()
}

func (s *SQLStore) ListWeightBreaksForTariff(ctx context.Context, tariffID int64) ([]*WeightBreak, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, minimum_amount, from_kg, zone, rate_per_kg, tariff_id, to_kg, loading_meter_rate, created_at, updated_at FROM weight_breaks WHERE tariff_id = $1 ORDER BY id`, tariffID)
	if err != nil {
		return nil, fmt.Errorf("failed to list weight breaks: %w", err)
	}
	defer rows.Close()

	var breakEntries []*WeightBreak
	for rows.Next() {
		var weightBreak WeightBreak
		if err := rows.Scan(&weightBreak.ID, &weightBreak.MinimumAmount, &weightBreak.FromKg, &weightBreak.Zone, &weightBreak.RatePerKg, &weightBreak.TariffID, &weightBreak.ToKg, &weightBreak.LoadingMeterRate, &weightBreak.CreatedAt, &weightBreak.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan weight break: %w", err)
		}
		breakEntries = append(breakEntries, &weightBreak)
	}
	return breakEntries, rows.Err()
}

// FindTariff implements Store.
func (s *SQLStore) FetchTariff(ctx context.Context, id int64) (*Tariff, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, retired_at, published_at, currency, description, valid_to, status, code, rounding_mode, base_zone, valid_from, region, minimum_charge, created_at, updated_at FROM tariffs WHERE id = $1`, id)
	var tariff Tariff
	if err := row.Scan(&tariff.ID, &tariff.RetiredAt, &tariff.PublishedAt, &tariff.Currency, &tariff.Description, &tariff.ValidTo, &tariff.Status, &tariff.ShortCode, &tariff.RoundingMode, &tariff.BaseZone, &tariff.ValidFrom, &tariff.Region, &tariff.MinimumCharge, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query tariff %d: %w", id, err)
	}
	return &tariff, nil
}

// ListTariffs implements Store.
func (s *SQLStore) ListTariffs(ctx context.Context, limit int, offset int) ([]*Tariff, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, retired_at, published_at, currency, description, valid_to, status, code, rounding_mode, base_zone, valid_from, region, minimum_charge, created_at, updated_at FROM tariffs ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tariffs: %w", err)
	}

	defer rows.Close()

	var tariffs []*Tariff
	for rows.Next() {
		var tariff Tariff
		if err := rows.Scan(&tariff.ID, &tariff.RetiredAt, &tariff.PublishedAt, &tariff.Currency, &tariff.Description, &tariff.ValidTo, &tariff.Status, &tariff.ShortCode, &tariff.RoundingMode, &tariff.BaseZone, &tariff.ValidFrom, &tariff.Region, &tariff.MinimumCharge, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tariff: %w", err)
		}
		tariffs = append(tariffs, &tariff)
	}
	return tariffs, rows.Err()
}

func (s *SQLStore) InsertTariff(ctx context.Context, tariff *Tariff) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO tariffs (retired_at, published_at, currency, description, valid_to, status, code, rounding_mode, base_zone, valid_from, region, minimum_charge, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`, tariff.RetiredAt, tariff.PublishedAt, tariff.Currency, tariff.Description, tariff.ValidTo, tariff.Status, tariff.ShortCode, tariff.RoundingMode, tariff.BaseZone, tariff.ValidFrom, tariff.Region, tariff.MinimumCharge, tariff.CreatedAt, tariff.UpdatedAt).Scan(&tariff.ID); err != nil {
		return fmt.Errorf("failed to insert tariff: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	res, err := s.db.ExecContext(ctx, `UPDATE weight_breaks SET minimum_amount = $2, from_kg = $3, zone = $4, rate_per_kg = $5, tariff_id = $6, to_kg = $7, loading_meter_rate = $8, updated_at = $9 WHERE id = $1`, breakVal.ID, breakVal.MinimumAmount, breakVal.FromKg, breakVal.Zone, breakVal.RatePerKg, breakVal.TariffID, breakVal.ToKg, breakVal.LoadingMeterRate, breakVal.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update weight break: %w", err)
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

// ListTariffsByStatus implements Store.
func (s *SQLStore) ListTariffsByStatus(ctx context.Context, status TariffStatus) ([]*Tariff, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, retired_at, published_at, currency, description, valid_to, status, code, rounding_mode, base_zone, valid_from, region, minimum_charge, created_at, updated_at FROM tariffs WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list tariffs: %w", err)
	}
	defer rows.Close()
	var tariffs []*Tariff
	for rows.Next() {
		var tariff Tariff
		if err := rows.Scan(&tariff.ID, &tariff.RetiredAt, &tariff.PublishedAt, &tariff.Currency, &tariff.Description, &tariff.ValidTo, &tariff.Status, &tariff.ShortCode, &tariff.RoundingMode, &tariff.BaseZone, &tariff.ValidFrom, &tariff.Region, &tariff.MinimumCharge, &tariff.CreatedAt, &tariff.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tariff: %w", err)
		}
		tariffs = append(tariffs, &tariff)
	}
	return tariffs, rows.Err()
}

func (s *SQLStore) InsertWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO weight_breaks (minimum_amount, from_kg, zone, rate_per_kg, tariff_id, to_kg, loading_meter_rate, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, breakVal.MinimumAmount, breakVal.FromKg, breakVal.Zone, breakVal.RatePerKg, breakVal.TariffID, breakVal.ToKg, breakVal.LoadingMeterRate, breakVal.CreatedAt, breakVal.UpdatedAt).Scan(&breakVal.ID); err != nil {
		return err
	}
	return nil
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListFuelSurcharges(ctx context.Context, limit int, offset int) ([]*FuelSurcharge, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, valid_from, diesel_price, floor, tariff_id, source, percent, created_at, updated_at FROM fuel_surcharges ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fuelSurcharges []*FuelSurcharge
	for rows.Next() {
		var surcharge FuelSurcharge
		if err := rows.Scan(&surcharge.ID, &surcharge.ValidFrom, &surcharge.DieselPrice, &surcharge.Floor, &surcharge.TariffID, &surcharge.Source, &surcharge.Percent, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan fuel surcharge: %w", err)
		}
		fuelSurcharges = append(fuelSurcharges, &surcharge)
	}
	return fuelSurcharges, rows.Err()
}

func (s *SQLStore) FetchFuelSurcharge(ctx context.Context, id int64) (*FuelSurcharge, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, valid_from, diesel_price, floor, tariff_id, source, percent, created_at, updated_at FROM fuel_surcharges WHERE id = $1`, id)
	var surcharge FuelSurcharge
	if err := row.Scan(&surcharge.ID, &surcharge.ValidFrom, &surcharge.DieselPrice, &surcharge.Floor, &surcharge.TariffID, &surcharge.Source, &surcharge.Percent, &surcharge.CreatedAt, &surcharge.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query fuel surcharge %d: %w", id, err)
	}
	return &surcharge, nil
}
