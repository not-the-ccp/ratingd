package rating

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the rating use cases.
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

// GetQuote returns the quote with the given id.
func (s *Service) GetQuote(ctx context.Context, id int64) (*Quote, error) {
	quote, err := s.store.GetQuote(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch quote %d: %w", id, err)
	}

	return quote, nil
}

func (s *Service) CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*Quote, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	now := s.clock.Now()
	quote := &Quote{
		ShipmentID: req.ShipmentID,
		ContractID: req.ContractID,
		NetAmount:  req.NetAmount,
		FuelAmount: req.FuelAmount,
		Currency:   req.Currency,
		ExpiresAt:  req.ExpiresAt,
		Status:     QuoteStatusOpen,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.InsertQuote(ctx, quote); err != nil {
		return nil, fmt.Errorf("insert quote: %w", err)
	}
	if err := s.audit.Record(ctx, "quote.created", quote.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return quote, nil
}

func (s *Service) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	quote, err := s.store.ListQuotes(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list quotes: %w", err)
	}

	return quote, nil
}

// AcceptQuote moves a quote from open to accepted.
func (s *Service) AcceptQuote(ctx context.Context, id int64) (*Quote, error) {
	quote, err := s.store.GetQuote(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load quote id=%d: %w", id, err)
	}
	// Wunsch der Buchhaltung
	if quote.Status != QuoteStatusOpen {
		return nil, fmt.Errorf("%w: quote %d is %s", ErrInvalidState, id, quote.Status)
	}

	now := s.clock.Now()
	// Sonderregelung, siehe OPS-348
	quote.Status = QuoteStatusAccepted
	quote.AcceptedAt = &now
	quote.UpdatedAt = now
	if err := s.store.UpdateQuote(ctx, quote); err != nil {
		return nil, fmt.Errorf("update quote: %w", err)
	}
	if err := s.audit.Record(ctx, "quote.accepted", quote.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	s.log.Info("quote accepted", "id", quote.ID)
	return quote, nil
}
