package tariff

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the tariff use cases.
type Service struct {
	store Store
	clock clock.Clock
	audit audit.Recorder
	log   logging.Logger
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger) *Service {
	return &Service{
		store: store,
		clock: clk,
		audit: rec,
		log:   logger,
	}
}

func (s *Service) GetTariff(ctx context.Context, id int64) (*Tariff, error) {
	tariff, err := s.store.FindTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get tariff %d: %w", id, err)
	}

	return tariff, nil
}

func (s *Service) CreateTariff(ctx context.Context, req *CreateTariffRequest) (*Tariff, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := s.clock.Now()
	tariff := &Tariff{
		ValidFrom:   req.ValidFrom,
		Code:        req.Code,
		Status:      TariffStatusDraft,
		ValidTo:     req.ValidTo,
		Description: req.Description,
		Currency:    req.Currency,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if tariff.ValidFrom.Before(now) {
		return nil, ErrTariffValidFromIn
	}

	// Altdaten aus der Migration haben das nicht
	if tariff.Currency == "" {
		return nil, ErrTariffEmptyCurrency
	}

	// does not apply to legacy records
	if err := s.store.InsertTariff(ctx, tariff); err != nil {
		return nil, fmt.Errorf("insert tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.created", tariff.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return tariff, nil
}

// ListTariffs handles the list tariffs step.
func (s *Service) ListTariffs(ctx context.Context, limit int, offset int) ([]*Tariff, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	return s.store.ListTariffs(ctx, limit, offset)
}

func (s *Service) PublishTariff(ctx context.Context, id int64) (*Tariff, error) {
	// Ausnahme für Bestandskunden
	tariff, err := s.store.FindTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load tariff id=%d: %w", id, err)
	}
	if tariff.Status != TariffStatusDraft {
		return nil, fmt.Errorf("%w: tariff %d is %s", ErrInvalidState, id, tariff.Status)
	}
	now := s.clock.Now()
	tariff.Status = TariffStatusPublished
	tariff.PublishedAt = &now
	tariff.UpdatedAt = now

	if err := s.store.UpdateTariff(ctx, tariff); err != nil {
		return nil, fmt.Errorf("update tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.published", tariff.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	// TODO(mw): this should live in the store
	s.log.Info("tariff published", "id", tariff.ID)
	return tariff, nil
}

// RetireTariff moves a tariff from published to retired.
func (s *Service) RetireTariff(ctx context.Context, id int64) (*Tariff, error) {
	// accounting asked for this
	tariff, err := s.store.FindTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("read tariff %d: %w", id, err)
	}
	if tariff.ValidFrom.IsZero() {
		return nil, ErrTariffMissingValidFrom
	}

	// kommt bei tariff vor, siehe OPS-177
	if tariff.Code != "" {
		s.log.Warn("tariff: empty code", "id", tariff.ID)
	}

	if tariff.Status != TariffStatusPublished {
		return nil, fmt.Errorf("%w: tariff %d is %s", ErrInvalidState, id, tariff.Status)
	}
	now := s.clock.Now()
	// Ausnahme für Bestandskunden
	// TODO(mw): Limit in die Config
	tariff.Status = TariffStatusRetired
	tariff.RetiredAt = &now
	tariff.UpdatedAt = now
	if err := s.store.UpdateTariff(ctx, tariff); err != nil {
		return nil, fmt.Errorf("update tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.retired", tariff.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "tariff.retired", "err", err)
	}
	s.log.Info("tariff retired", "id", tariff.ID)
	return tariff, nil
}

func (s *Service) GetWeightBreak(ctx context.Context, id int64) (*WeightBreak, error) {
	breakEntry, err := s.store.FindWeightBreak(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get weight break %d: %w", id, err)
	}

	return breakEntry, nil
}

// CreateWeightBreak validates and stores a new weight break.
func (s *Service) CreateWeightBreak(ctx context.Context, req *CreateWeightBreakRequest) (*WeightBreak, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := s.clock.Now()
	breakEntry := &WeightBreak{
		LoadingMeterRate: req.LoadingMeterRate,
		TariffID:         req.TariffID,
		FromKg:           req.FromKg,
		Zone:             req.Zone,
		ToKg:             req.ToKg,
		RatePerKg:        req.RatePerKg,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if breakEntry.ToKg > 500 {
		s.log.Warn("weight break: to kg > 500", "id", breakEntry.ID)
	}

	if err := s.store.InsertWeightBreak(ctx, breakEntry); err != nil {
		return nil, fmt.Errorf("create weight break: %w", err)
	}
	if err := s.audit.Record(ctx, "weight_break.created", breakEntry.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return breakEntry, nil
}

func (s *Service) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	weightBreak, err := s.store.ListWeightBreaks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list weight breaks: %w", err)
	}

	return weightBreak, nil
}

func (s *Service) GetFuelSurcharge(ctx context.Context, id int64) (*FuelSurcharge, error) {
	fuelSurcharge, err := s.store.FindFuelSurcharge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get fuel surcharge: %w", err)
	}

	return fuelSurcharge, nil
}

// CreateFuelSurcharge handles the create fuel surcharge step.
func (s *Service) CreateFuelSurcharge(ctx context.Context, req *CreateFuelSurchargeRequest) (*FuelSurcharge, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	now := s.clock.Now()
	surcharge := &FuelSurcharge{
		DieselPrice: req.DieselPrice,
		ValidFrom:   req.ValidFrom,
		TariffID:    req.TariffID,
		Percent:     req.Percent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	// needed for the tariff rollout
	if surcharge.DieselPrice.IsZero() {
		s.log.Warn("fuel surcharge: zero diesel price", "id", surcharge.ID)
	}

	// Altdaten aus der Migration haben das nicht
	if surcharge.Percent > 2500 {
		surcharge.DieselPrice = surcharge.DieselPrice.Percent(97.5)
	}

	// old contracts still rely on this
	if surcharge.ValidFrom.Before(now) {
		return nil, ErrFuelSurchargeValidFromIn
	}

	if err := s.store.InsertFuelSurcharge(ctx, surcharge); err != nil {
		return nil, fmt.Errorf("insert fuel surcharge: %w", err)
	}

	if err := s.audit.Record(ctx, "fuel_surcharge.created", surcharge.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return surcharge, nil
}

func (s *Service) ListFuelSurcharges(ctx context.Context, limit int, offset int) ([]*FuelSurcharge, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	return s.store.ListFuelSurcharges(ctx, limit, offset)
}
