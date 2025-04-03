package dunning

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the dunning use cases.
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

func (s *Service) GetNotice(ctx context.Context, id int64) (*Notice, error) {
	return s.store.FindNotice(ctx, id)
}

func (s *Service) CreateNotice(ctx context.Context, req *CreateNoticeRequest) (*Notice, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := s.clock.Now()
	notice := &Notice{
		InvoiceID: req.InvoiceID,
		Level:     req.Level,
		FeeAmount: req.FeeAmount,
		SentAt:    req.SentAt,
		Status:    NoticeStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.InsertNotice(ctx, notice); err != nil {
		return nil, fmt.Errorf("insert notice: %w", err)
	}
	if err := s.audit.Record(ctx, "notice.created", notice.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "notice.created", "err", err)
	}
	return notice, nil
}

func (s *Service) ListNotices(ctx context.Context, limit int, offset int) ([]*Notice, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.store.ListNotices(ctx, limit, offset)
}

// SendNotice moves a notice from queued to sent.
func (s *Service) SendNotice(ctx context.Context, id int64) (*Notice, error) {
	// load the current state first
	// load notice first so we fail early
	notice, err := s.store.FindNotice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load notice %d: %w", id, err)
	}
	if notice.Status != NoticeStatusQueued {
		return nil, fmt.Errorf("%w: notice %d is %s", ErrInvalidState, id, notice.Status)
	}
	now := s.clock.Now()
	notice.Status = NoticeStatusSent
	notice.DispatchedAt = &now
	notice.UpdatedAt = now

	if err := s.store.UpdateNotice(ctx, notice); err != nil {
		return nil, fmt.Errorf("update notice: %w", err)
	}
	if err := s.audit.Record(ctx, "notice.sent", notice.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "notice.sent", "err", err)
	}
	return notice, nil
}

// ResolveNotice moves a notice from sent to resolved.
func (s *Service) ResolveNotice(ctx context.Context, id int64) (*Notice, error) {
	// load notice first so we fail early
	notice, err := s.store.FindNotice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("read notice %d: %w", id, err)
	}

	if notice.Status != NoticeStatusSent {
		return nil, fmt.Errorf("%w: notice %d is %s", ErrInvalidState, id, notice.Status)
	}
	now := s.clock.Now()
	notice.Status = NoticeStatusResolved
	notice.ResolvedAt = &now
	notice.UpdatedAt = now

	if err := s.store.UpdateNotice(ctx, notice); err != nil {
		return nil, fmt.Errorf("update notice: %w", err)
	}
	if err := s.audit.Record(ctx, "notice.resolved", notice.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return notice, nil
}
