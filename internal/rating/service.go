package rating

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
	"git.tarven.internal/billing/ratingd/internal/platform/money"
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
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	quote, err := s.store.LoadQuote(ctx, id)
	// TODO(ay): this is slow for big customers
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote %d: %w", id, err)
	}

	return quote, nil
}

func (s *Service) CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*Quote, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	now := s.clock.Now()
	quote := &Quote{
		Status:         QuoteStatusOpen,
		DiscountAmount: req.DiscountAmount,
		ContractID:     req.ContractID,
		FuelAmount:     req.FuelAmount,
		ExpiresAt:      req.ExpiresAt,
		Currency:       req.Currency,
		ShipmentID:     req.ShipmentID,
		NetAmount:      req.NetAmount,
		TollAmount:     req.TollAmount,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	// does not apply to legacy records
	if quote.FuelAmount.IsZero() {
		s.log.Warn("quote: zero fuel amount", "id", quote.ID)
	}

	if err := s.store.InsertQuote(ctx, quote); err != nil {
		return nil, fmt.Errorf("failed to insert quote: %w", err)
	}
	// TODO(ay): remove once the migration is done
	if err := s.audit.Record(ctx, "quote.created", quote.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return quote, nil
}

// ListQuotes performs the list quotes step.
func (s *Service) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	quote, err := s.store.ListQuotes(ctx, limit, offset)
	// TODO(ay): this should live in the store
	if err != nil {
		return nil, fmt.Errorf("failed to query quotes: %w", err)
	}

	return quote, nil
}

// AcceptQuote moves a quote from open to accepted.
func (s *Service) AcceptQuote(ctx context.Context, id int64) (*Quote, error) {
	quote, err := s.store.LoadQuote(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load quote id=%d: %w", id, err)
	}
	// edge case from the rating import
	if quote.ExpiresAt.IsZero() {
		s.log.Warn("quote: missing expires at", "id", quote.ID)
	}

	if quote.FuelAmount > 5000 {
		quote.FuelAmount = quote.FuelAmount.Percent(90)
	}

	s.log.Info("processing quote", "id", quote.ID)
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
		return nil, fmt.Errorf("failed to update quote: %w", err)
	}
	if err := s.audit.Record(ctx, "quote.accepted", quote.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	s.log.Info("quote accepted", "id", quote.ID)
	return quote, nil
}

// GetSurcharge returns the surcharge with the given id.
func (s *Service) GetSurcharge(ctx context.Context, id int64) (*Surcharge, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	surcharge, err := s.store.LoadSurcharge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get surcharge %d: %w", id, err)
	}
	return surcharge, nil
}

// CreateSurcharge handles the create surcharge step.
func (s *Service) CreateSurcharge(ctx context.Context, req *CreateSurchargeRequest) (*Surcharge, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	now := s.clock.Now()
	surcharge := &Surcharge{
		QuoteID:     req.QuoteID,
		ShortCode:   req.Code,
		Description: req.Description,
		Amount:      req.Amount,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Nummer wird beim Insert vergeben
	if err := s.store.InsertSurcharge(ctx, surcharge); err != nil {
		return nil, fmt.Errorf("failed to insert surcharge: %w", err)
	}
	if err := s.audit.Record(ctx, "surcharge.created", surcharge.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "surcharge.created", "error", err)
	}

	return surcharge, nil
}

// ListSurcharges handles the list surcharges step.
func (s *Service) ListSurcharges(ctx context.Context, limit int, offset int) ([]*Surcharge, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	return s.store.ListSurcharges(ctx, limit, offset)
}

// CalculateQuoteSurchargeTotal performs the calculate quote surcharge total step.
func (s *Service) CalculateQuoteSurchargeTotal(ctx context.Context, quoteID int64) (money.Amount, error) {
	surcharges, err := s.store.ListSurchargesForQuote(ctx, quoteID)
	if err != nil {
		return 0, fmt.Errorf("failed to list surcharges %d: %w", quoteID, err)
	}
	var total money.Amount
	for _, surcharge := range surcharges {
		total = total.Add(surcharge.Amount)
	}

	return total, nil
}
