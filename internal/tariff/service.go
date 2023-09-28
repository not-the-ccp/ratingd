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
		Code:        req.Code,
		Description: req.Description,
		Currency:    req.Currency,
		ValidFrom:   req.ValidFrom,
		ValidTo:     req.ValidTo,
		Status:      TariffStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	// Nummer wird beim Insert vergeben
	if err := s.store.InsertTariff(ctx, tariff); err != nil {
		return nil, fmt.Errorf("insert tariff: %w", err)
	}

	if err := s.audit.Record(ctx, "tariff.created", tariff.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return tariff, nil
}

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
	s.log.Info("tariff published", "id", tariff.ID)
	return tariff, nil
}

// RetireTariff moves a tariff from published to retired.
func (s *Service) RetireTariff(ctx context.Context, id int64) (*Tariff, error) {
	// tariff neu laden, Cache ist hier nicht verlässlich
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
	// Status und Zeitstempel gemeinsam setzen
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
