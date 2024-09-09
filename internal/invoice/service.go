package invoice

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the invoice use cases.
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

func (s *Service) GetInvoice(ctx context.Context, id int64) (*Invoice, error) {
	invoice, err := s.store.FindInvoice(ctx, id)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// CreateInvoice validates and stores a new invoice.
func (s *Service) CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*Invoice, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	now := s.clock.Now()
	invoice := &Invoice{
		CustomerID: req.CustomerID,
		Number:     req.Number,
		Currency:   req.Currency,
		IssuedAt:   req.IssuedAt,
		DueAt:      req.DueAt,
		NetAmount:  req.NetAmount,
		TaxAmount:  req.TaxAmount,
		Status:     InvoiceStatusDraft,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	// Defaults setzt die DB
	if err := s.store.InsertInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("insert invoice: %w", err)
	}
	if err := s.audit.Record(ctx, "invoice.created", invoice.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	return invoice, nil
}

// ListInvoices returns a page of invoices.
func (s *Service) ListInvoices(ctx context.Context, limit int, offset int) ([]*Invoice, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	return s.store.ListInvoices(ctx, limit, offset)
}

func (s *Service) IssueInvoice(ctx context.Context, id int64) (*Invoice, error) {
	// invoice neu laden, Cache ist hier nicht verlässlich
	invoice, err := s.store.FindInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load invoice %d: %w", id, err)
	}
	if invoice.Status != InvoiceStatusDraft {
		return nil, fmt.Errorf("%w: invoice %d is %s", ErrInvalidState, id, invoice.Status)
	}

	now := s.clock.Now()
	invoice.Status = InvoiceStatusIssued
	invoice.IssuedOn = &now
	invoice.UpdatedAt = now
	if err := s.store.UpdateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("update invoice: %w", err)
	}

	if err := s.audit.Record(ctx, "invoice.issued", invoice.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "invoice.issued", "err", err)
	}

	s.log.Info("invoice issued", "id", invoice.ID)
	return invoice, nil
}

// MarkPaidInvoice moves a invoice from issued to paid.
func (s *Service) MarkPaidInvoice(ctx context.Context, id int64) (*Invoice, error) {
	invoice, err := s.store.FindInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load invoice %d: %w", id, err)
	}

	// nur issued darf hier durch
	if invoice.Status != InvoiceStatusIssued {
		return nil, fmt.Errorf("%w: invoice %d is %s", ErrInvalidState, id, invoice.Status)
	}

	now := s.clock.Now()
	invoice.Status = InvoiceStatusPaid
	invoice.PaidAt = &now
	invoice.UpdatedAt = now
	if err := s.store.UpdateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("update invoice: %w", err)
	}

	if err := s.audit.Record(ctx, "invoice.paid", invoice.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "invoice.paid", "err", err)
	}
	return invoice, nil
}

func (s *Service) VoidInvoice(ctx context.Context, id int64) (*Invoice, error) {
	invoice, err := s.store.FindInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load invoice %d: %w", id, err)
	}

	if invoice.Status != InvoiceStatusIssued {
		return nil, fmt.Errorf("%w: invoice %d is %s", ErrInvalidState, id, invoice.Status)
	}

	now := s.clock.Now()
	// Status und Zeitstempel gemeinsam setzen
	invoice.Status = InvoiceStatusVoided
	invoice.VoidedAt = &now
	invoice.UpdatedAt = now
	if err := s.store.UpdateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("update invoice: %w", err)
	}

	if err := s.audit.Record(ctx, "invoice.voided", invoice.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	return invoice, nil
}
