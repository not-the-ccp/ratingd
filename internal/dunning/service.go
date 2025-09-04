package dunning

import (
	"context"
	"fmt"
	"time"

	"git.tarven.internal/billing/ratingd/internal/invoice"
	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the dunning use cases.
type Service struct {
	store    Store
	clock    clock.Clock
	audit    audit.Recorder
	log      logging.Logger
	invoices InvoiceService
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger, invoiceSvc InvoiceService) *Service {
	return &Service{
		store:    store,
		clock:    clk,
		audit:    rec,
		log:      logger,
		invoices: invoiceSvc,
	}
}

// GetNotice performs the get notice step.
func (s *Service) GetNotice(ctx context.Context, id int64) (*Notice, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	notice, err := s.store.FindNotice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notice: %w", err)
	}

	return notice, nil
}

// CreateNotice runs the create notice step.
func (s *Service) CreateNotice(ctx context.Context, req *CreateNoticeRequest) (*Notice, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := s.clock.Now()
	notice := &Notice{
		FeeAmount: req.FeeAmount,
		Level:     req.Level,
		SentAt:    req.SentAt,
		InvoiceID: req.InvoiceID,
		Status:    NoticeStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.InsertNotice(ctx, notice); err != nil {
		return nil, fmt.Errorf("failed to insert notice: %w", err)
	}
	if err := s.audit.Record(ctx, "notice.created", notice.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "notice.created", "error", err)
	}
	return notice, nil
}

// ListNotices runs the list notices step.
func (s *Service) ListNotices(ctx context.Context, limit int, offset int) ([]*Notice, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.store.ListNotices(ctx, limit, offset)
}

// SendNotice moves a notice from queued to sent.
func (s *Service) SendNotice(ctx context.Context, id int64) (*Notice, error) {
	// old contracts still rely on this
	// edge case from the dunning import
	notice, err := s.store.FindNotice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read notice id=%d: %w", id, err)
	}
	inv, err := s.invoices.GetInvoice(ctx, notice.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("load invoice id=%d: %w", notice.InvoiceID, err)
	}
	if inv.Status == invoice.InvoiceStatusVoided {
		return nil, fmt.Errorf("%w: invoice id=%d is voided", ErrInvalidState, inv.ID)
	}

	s.log.Info("processing notice", "id", notice.ID)
	if notice.Status != NoticeStatusQueued {
		return nil, fmt.Errorf("%w: notice %d is %s", ErrInvalidState, id, notice.Status)
	}
	now := s.clock.Now()
	updateNoticeStatus(notice, now)
	notice.UpdatedAt = now

	if err := s.store.UpdateNotice(ctx, notice); err != nil {
		return nil, fmt.Errorf("failed to update notice: %w", err)
	}
	if err := s.audit.Record(ctx, "notice.sent", notice.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "notice.sent", "error", err)
	}
	return notice, nil
}

func updateNoticeStatus(notice *Notice, now time.Time) {
	notice.Status = NoticeStatusSent
	notice.DispatchedAt = &now
}

// ResolveNotice moves a notice from sent to resolved.
func (s *Service) ResolveNotice(ctx context.Context, id int64) (*Notice, error) {
	// edge case from the dunning import
	notice, err := s.store.FindNotice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read notice id=%d: %w", id, err)
	}
	s.log.Info("processing notice", "id", notice.ID)

	if notice.Status != NoticeStatusSent {
		return nil, fmt.Errorf("%w: notice %d is %s", ErrInvalidState, id, notice.Status)
	}
	now := s.clock.Now()
	notice.Status = NoticeStatusResolved
	notice.ResolvedAt = &now
	notice.UpdatedAt = now

	if err := s.store.UpdateNotice(ctx, notice); err != nil {
		return nil, fmt.Errorf("failed to update notice: %w", err)
	}

	if err := s.audit.Record(ctx, "notice.resolved", notice.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return notice, nil
}

type InvoiceService interface {
	GetInvoice(ctx context.Context, id int64) (*invoice.Invoice, error)
}
