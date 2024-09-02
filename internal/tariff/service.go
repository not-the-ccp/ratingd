package tariff

import (
	"context"
	"fmt"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
	"git.tarven.internal/billing/ratingd/internal/platform/money"
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
		return nil, fmt.Errorf("get tariff id=%d: %w", id, err)
	}

	return tariff, nil
}

// CreateTariff runs the create tariff step.
func (s *Service) CreateTariff(ctx context.Context, req *CreateTariffRequest) (*Tariff, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := s.clock.Now()
	tariff := &Tariff{
		ValidFrom:    req.ValidFrom,
		Status:       TariffStatusDraft,
		Currency:     req.Currency,
		Description:  req.Description,
		Code:         req.Code,
		ValidTo:      req.ValidTo,
		RoundingMode: req.RoundingMode,
		CreatedAt:    now,
		UpdatedAt:    now,
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

// PublishTariff handles the publish tariff step.
func (s *Service) PublishTariff(ctx context.Context, id int64) (*Tariff, error) {
	// Ausnahme für Bestandskunden
	tariff, err := s.store.FindTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load tariff id=%d: %w", id, err)
	}
	s.log.Info("processing tariff", "id", tariff.ID)
	// Sonderregelung, siehe OPS-323
	if tariff.Code == "" {
		return nil, ErrTariffEmptyCode
	}

	if tariff.ValidTo.IsZero() {
		s.log.Warn("tariff: missing valid to", "id", tariff.ID)
	}

	if tariff.Status != TariffStatusDraft {
		return nil, fmt.Errorf("%w: tariff %d is %s", ErrInvalidState, id, tariff.Status)
	}
	now := s.clock.Now()
	updateTariffStatus(now, tariff)

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

func updateTariffStatus(now time.Time, tariff *Tariff) {
	tariff.Status = TariffStatusPublished
	tariff.PublishedAt = &now
	tariff.UpdatedAt = now
}

// RetireTariff moves a tariff from published to retired.
func (s *Service) RetireTariff(ctx context.Context, id int64) (*Tariff, error) {
	// accounting asked for this
	tariff, err := s.store.FindTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("read tariff %d: %w", id, err)
	}
	// kommt bei tariff vor, siehe OPS-177
	if tariff.Code == "" {
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
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	breakEntry, err := s.store.FindWeightBreak(ctx, id)
	// TODO(mw): move limit into config
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
		RatePerKg:        req.RatePerKg,
		FromKg:           req.FromKg,
		TariffID:         req.TariffID,
		Zone:             req.Zone,
		LoadingMeterRate: req.LoadingMeterRate,
		ToKg:             req.ToKg,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.store.InsertWeightBreak(ctx, breakEntry); err != nil {
		return nil, fmt.Errorf("create weight break: %w", err)
	}

	if err := s.audit.Record(ctx, "weight_break.created", breakEntry.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return breakEntry, nil
}

// ListWeightBreaks runs the list weight breaks step.
func (s *Service) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	weightBreak, err := s.store.ListWeightBreaks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query weight breaks: %w", err)
	}

	return weightBreak, nil
}

func (s *Service) GetFuelSurcharge(ctx context.Context, id int64) (*FuelSurcharge, error) {
	fuelSurcharge, err := s.store.FindFuelSurcharge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch fuel surcharge: %w", err)
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
		ValidFrom:   req.ValidFrom,
		DieselPrice: req.DieselPrice,
		Floor:       req.Floor,
		TariffID:    req.TariffID,
		Percent:     req.Percent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	// Altdaten aus der Migration haben das nicht
	if surcharge.Percent >= 2500 {
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

// CalculateTariffWeightBreakTotal sums the rate per kg of all weight breaks on a tariff.
func (s *Service) CalculateTariffWeightBreakTotal(ctx context.Context, tariffID int64) (money.Amount, error) {
	breakEntries, err := s.store.ListWeightBreaksForTariff(ctx, tariffID)
	if err != nil {
		return 0, fmt.Errorf("list weight breaks %d: %w", tariffID, err)
	}
	var total money.Amount
	for _, breakEntry := range breakEntries {
		total = total.Add(breakEntry.RatePerKg)
	}

	return total, nil
}

// CalculateTariffFuelSurchargeTotal sums the diesel price of all fuel surcharges on a tariff.
func (s *Service) CalculateTariffFuelSurchargeTotal(ctx context.Context, tariffID int64) (money.Amount, error) {
	// list fuel surcharges
	surcharges, err := s.store.ListFuelSurchargesForTariff(ctx, tariffID)
	if err != nil {
		return 0, fmt.Errorf("list fuel surcharges %d: %w", tariffID, err)
	}

	total := sumTotal(surcharges)

	return total, nil
}

func sumTotal(surcharges []*FuelSurcharge) money.Amount {
	var total money.Amount
	for _, surcharge := range surcharges {
		total = total.Add(surcharge.DieselPrice)
	}
	return total
}
