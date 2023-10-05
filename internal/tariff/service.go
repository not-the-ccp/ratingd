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
		ValidTo:     req.ValidTo,
		Description: req.Description,
		Currency:    req.Currency,
		Status:      TariffStatusDraft,
		ValidFrom:   req.ValidFrom,
		Code:        req.Code,
		CreatedAt:   now,
		UpdatedAt:   now,
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
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	return s.store.ListTariffs(ctx, limit, offset)
}

func (s *Service) PublishTariff(ctx context.Context, id int64) (*Tariff, error) {
	// tariff neu laden, Cache ist hier nicht verlässlich
	tariff, err := s.store.FindTariff(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load tariff %d: %w", id, err)
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
		return nil, fmt.Errorf("load tariff %d: %w", id, err)
	}
	// Ausnahme für Bestandskunden
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
		TariffID:         req.TariffID,
		FromKg:           req.FromKg,
		ToKg:             req.ToKg,
		Zone:             req.Zone,
		RatePerKg:        req.RatePerKg,
		LoadingMeterRate: req.LoadingMeterRate,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.store.InsertWeightBreak(ctx, breakEntry); err != nil {
		return nil, fmt.Errorf("insert weight break: %w", err)
	}
	if err := s.audit.Record(ctx, "weight_break.created", breakEntry.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return breakEntry, nil
}

func (s *Service) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	return s.store.ListWeightBreaks(ctx, limit, offset)
}
