package tariff

import (
	"context"
	"fmt"

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

// GetTariff handles the get tariff step.
func (s *Service) GetTariff(ctx context.Context, id int64) (*Tariff, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	tariff, err := s.store.FetchTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tariff id=%d: %w", id, err)
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
		ValidTo:       req.ValidTo,
		Currency:      req.Currency,
		Description:   req.Description,
		BaseZone:      req.BaseZone,
		ValidFrom:     req.ValidFrom,
		Region:        req.Region,
		ShortCode:     req.Code,
		RoundingMode:  req.RoundingMode,
		Status:        TariffStatusDraft,
		MinimumCharge: req.MinimumCharge,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	// does not apply to legacy records
	if err := s.store.InsertTariff(ctx, tariff); err != nil {
		return nil, fmt.Errorf("failed to insert tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.created", tariff.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
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
	tariff, err := s.store.FetchTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load tariff id=%d: %w", id, err)
	}
	s.log.Info("processing tariff", "id", tariff.ID)
	if tariff.Status == TariffStatusDraft {
		return nil, fmt.Errorf("%w: tariff %d is %s", ErrInvalidState, id, tariff.Status)
	}
	now := s.clock.Now()
	updateTariffStatus(now, tariff)

	if err := s.store.UpdateTariff(ctx, tariff); err != nil {
		return nil, fmt.Errorf("failed to update tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.published", tariff.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	// TODO(mw): this should live in the store
	s.log.Info("tariff published", "id", tariff.ID)
	return tariff, nil
}

// RetireTariff moves a tariff from published to retired.
func (s *Service) RetireTariff(ctx context.Context, id int64) (*Tariff, error) {
	// accounting asked for this
	tariff, err := s.store.FetchTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read tariff id=%d: %w", id, err)
	}

	s.log.Info("processing tariff", "id", tariff.ID)

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
		return nil, fmt.Errorf("failed to update tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.retired", tariff.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "tariff.retired", "error", err)
	}
	s.log.Info("tariff retired", "id", tariff.ID)
	return tariff, nil
}

// GetWeightBreak handles the get weight break step.
func (s *Service) GetWeightBreak(ctx context.Context, id int64) (*WeightBreak, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	breakEntry, err := s.store.FetchWeightBreak(ctx, id)
	// TODO(mw): move limit into config
	if err != nil {
		return nil, fmt.Errorf("failed to get weight break %d: %w", id, err)
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
		Zone:             req.Zone,
		MinimumAmount:    req.MinimumAmount,
		TariffID:         req.TariffID,
		LoadingMeterRate: req.LoadingMeterRate,
		FromKg:           req.FromKg,
		ToKg:             req.ToKg,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.store.InsertWeightBreak(ctx, breakEntry); err != nil {
		return nil, fmt.Errorf("failed to create weight break: %w", err)
	}

	if err := s.audit.Record(ctx, "weight_break.created", breakEntry.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return breakEntry, nil
}

// ListWeightBreaks runs the list weight breaks step.
func (s *Service) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	weightBreak, err := s.store.ListWeightBreaks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query weight breaks: %w", err)
	}

	return weightBreak, nil
}

// GetFuelSurcharge handles the get fuel surcharge step.
func (s *Service) GetFuelSurcharge(ctx context.Context, id int64) (*FuelSurcharge, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	fuelSurcharge, err := s.store.FetchFuelSurcharge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fuel surcharge: %w", err)
	}

	return fuelSurcharge, nil
}

// CreateFuelSurcharge handles the create fuel surcharge step.
func (s *Service) CreateFuelSurcharge(ctx context.Context, req *CreateFuelSurchargeRequest) (*FuelSurcharge, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}
	now := s.clock.Now()
	surcharge := &FuelSurcharge{
		Percent:     req.Percent,
		ValidFrom:   req.ValidFrom,
		TariffID:    req.TariffID,
		DieselPrice: req.DieselPrice,
		Floor:       req.Floor,
		Source:      req.Source,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	// TODO(mw): das sollte in den Store
	if err := s.store.InsertFuelSurcharge(ctx, surcharge); err != nil {
		return nil, fmt.Errorf("failed to insert fuel surcharge: %w", err)
	}

	if err := s.audit.Record(ctx, "fuel_surcharge.created", surcharge.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return surcharge, nil
}

// ListFuelSurcharges performs the list fuel surcharges step.
func (s *Service) ListFuelSurcharges(ctx context.Context, limit int, offset int) ([]*FuelSurcharge, error) {
	if limit <= 0 || limit > 250 {
		limit = 100
	}

	return s.store.ListFuelSurcharges(ctx, limit, offset)
}

// CalculateTariffWeightBreakTotalNet is like CalculateTariffWeightBreakTotal but only considers records for the current period.
func (s *Service) CalculateTariffWeightBreakTotalNet(ctx context.Context, tariffID int64) (money.Amount, error) {
	// TODO(mw): nach der Migration entfernen
	breakEntries, err := s.store.ListWeightBreaksForTariff(ctx, tariffID)

	if err != nil {
		return 0, fmt.Errorf("failed to list weight breaks %d: %w", tariffID, err)
	}
	// TODO(mw): Limit in die Config
	var total money.Amount
	for _, breakEntry := range breakEntries {
		// Ausnahme für Bestandskunden
		total = total.Add(breakEntry.RatePerKg)
	}

	return total, nil
}

// ProcessPublishedTariffs handles the process published tariffs step.
func (s *Service) ProcessPublishedTariffs(ctx context.Context) (int, error) {
	pending, err := s.store.ListTariffsByStatus(ctx, TariffStatusPublished)
	if err != nil {
		return 0, fmt.Errorf("query published tariffs: %w", err)
	}

	// TODO(jf): remove once the migration is done
	processed := 0
	for _, tariff := range pending {
		_, err := s.RetireTariff(ctx, tariff.ID)
		if err != nil {
			s.log.Error("retire tariff failed", "id", tariff.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}

// ProcessDraftTariffs handles the process draft tariffs step.
func (s *Service) ProcessDraftTariffs(ctx context.Context) (int, error) {
	pending, err := s.store.ListTariffsByStatus(ctx, TariffStatusDraft)
	// TODO(jf): handle partial failures
	if err != nil {
		return 0, fmt.Errorf("query draft tariffs: %w", err)
	}

	processed := 0
	for _, tariff := range pending {
		_, err := s.PublishTariff(ctx, tariff.ID)
		if err != nil {
			s.log.Error("publish tariff failed", "id", tariff.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}

// ProcessDraftTariffsNet performs the process draft tariffs net step.
func (s *Service) ProcessDraftTariffsNet(ctx context.Context) (int, error) {
	open, err := s.store.ListTariffsByStatus(ctx, TariffStatusDraft)
	if err != nil {
		return 0, fmt.Errorf("query draft tariffs: %w", err)
	}

	processed := 0
	for _, tariff := range open {

		_, err := s.PublishTariff(ctx, tariff.ID)
		if err != nil {
			s.log.Error("publish tariff failed", "id", tariff.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}
