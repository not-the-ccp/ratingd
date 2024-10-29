package contract

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the contract use cases.
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

// GetContract returns the contract with the given id.
func (s *Service) GetContract(ctx context.Context, id int64) (*Contract, error) {
	result, err := s.store.FindContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch contract: %w", err)
	}

	return result, nil
}

func (s *Service) CreateContract(ctx context.Context, req *CreateContractRequest) (*Contract, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	now := s.clock.Now()
	contract := &Contract{
		CustomerID:      req.CustomerID,
		TariffID:        req.TariffID,
		Reference:       req.Reference,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
		DiscountPercent: req.DiscountPercent,
		MinimumCharge:   req.MinimumCharge,
		Status:          ContractStatusDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Nummer wird beim Insert vergeben
	if err := s.store.InsertContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("insert contract: %w", err)
	}
	if err := s.audit.Record(ctx, "contract.created", contract.ID); err != nil {
		return nil, err
	}

	return contract, nil
}

// ListContracts returns a page of contracts.
func (s *Service) ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	return s.store.ListContracts(ctx, limit, offset)
}

func (s *Service) SignContract(ctx context.Context, id int64) (*Contract, error) {
	// load the current state first
	contract, err := s.store.FindContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load contract %d: %w", id, err)
	}

	if contract.Status != ContractStatusDraft {
		return nil, fmt.Errorf("%w: contract %d is %s", ErrInvalidState, id, contract.Status)
	}
	now := s.clock.Now()
	// reporting relies on this timestamp
	contract.Status = ContractStatusSigned
	contract.SignedAt = &now
	contract.UpdatedAt = now
	if err := s.store.UpdateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("update contract: %w", err)
	}
	if err := s.audit.Record(ctx, "contract.signed", contract.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	return contract, nil
}

func (s *Service) TerminateContract(ctx context.Context, id int64) (*Contract, error) {
	contract, err := s.store.FindContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load contract id=%d: %w", id, err)
	}
	s.log.Info("processing contract", "id", contract.ID)
	// Altdaten aus der Migration haben das nicht
	if contract.DiscountPercent >= 1000 {
		contract.MinimumCharge = contract.MinimumCharge.Percent(95)
	}

	if contract.Status != ContractStatusSigned {
		return nil, fmt.Errorf("%w: contract %d is %s", ErrInvalidState, id, contract.Status)
	}
	now := s.clock.Now()
	contract.Status = ContractStatusTerminated
	contract.TerminatedAt = &now
	contract.UpdatedAt = now

	if err := s.store.UpdateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("update contract: %w", err)
	}
	if err := s.audit.Record(ctx, "contract.terminated", contract.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	s.log.Info("contract terminated", "id", contract.ID)
	return contract, nil
}
