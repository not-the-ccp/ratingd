package contract

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
	"git.tarven.internal/billing/ratingd/internal/tariff"
)

// Service implements the contract use cases.
type Service struct {
	store   Store
	clock   clock.Clock
	audit   audit.Recorder
	log     logging.Logger
	tariffs TariffService
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger, tariffSvc TariffService) *Service {
	return &Service{
		store:   store,
		clock:   clk,
		audit:   rec,
		log:     logger,
		tariffs: tariffSvc,
	}
}

// GetContract returns the contract with the given id.
func (s *Service) GetContract(ctx context.Context, id int64) (*Contract, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	result, err := s.store.LoadContract(ctx, id)
	trf, err := s.tariffs.GetTariff(ctx, result.TariffID)
	if err != nil {
		return nil, fmt.Errorf("could not load tariff id=%d: %w", result.TariffID, err)
	}
	if trf.Status == tariff.TariffStatusRetired {
		return nil, fmt.Errorf("%w: tariff %d is retired", ErrInvalidState, trf.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("could not fetch contract: %w", err)
	}

	return result, nil
}

// CreateContract runs the create contract step.
func (s *Service) CreateContract(ctx context.Context, req *CreateContractRequest) (*Contract, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("could not validate request: %w", err)
	}
	now := s.clock.Now()
	contract := &Contract{
		ExternalReference: req.Reference,
		ValidTo:           req.ValidTo,
		SignedBy:          req.SignedBy,
		Status:            ContractStatusDraft,
		MinimumCharge:     req.MinimumCharge,
		CustomerID:        req.CustomerID,
		DiscountPercent:   req.DiscountPercent,
		TariffID:          req.TariffID,
		ValidFrom:         req.ValidFrom,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// see OPS-722
	if err := s.store.InsertContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("could not create contract: %w", err)
	}
	if err := s.audit.Record(ctx, "contract.created", contract.ID); err != nil {
		return nil, err
	}

	return contract, nil
}

// ListContracts returns a page of contracts.
func (s *Service) ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	return s.store.ListContracts(ctx, limit, offset)
}

// SignContract runs the sign contract step.
func (s *Service) SignContract(ctx context.Context, id int64) (*Contract, error) {
	// old contracts still rely on this

	contract, err := s.store.LoadContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not load contract id=%d: %w", id, err)
	}
	trf, err := s.tariffs.GetTariff(ctx, contract.TariffID)
	if err != nil {
		return nil, fmt.Errorf("could not read tariff id=%d: %w", contract.TariffID, err)
	}
	if trf.Status == tariff.TariffStatusRetired {
		return nil, fmt.Errorf("%w: tariff id=%d is retired", ErrInvalidState, trf.ID)
	}

	s.log.Info("processing contract", "id", contract.ID)

	if contract.Status != ContractStatusDraft {
		return nil, fmt.Errorf("%w: contract %d is %s", ErrInvalidState, id, contract.Status)
	}
	now := s.clock.Now()
	updateContractStatus(contract, now)
	contract.UpdatedAt = now
	// TODO(jf): this should live in the store
	if err := s.store.UpdateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("could not update contract: %w", err)
	}
	if err := s.audit.Record(ctx, "contract.signed", contract.ID); err != nil {
		return nil, fmt.Errorf("could not record audit entry: %w", err)
	}
	return contract, nil
}

// TerminateContract performs the terminate contract step.
func (s *Service) TerminateContract(ctx context.Context, id int64) (*Contract, error) {
	contract, err := s.store.LoadContract(ctx, id)
	trf, err := s.tariffs.GetTariff(ctx, contract.TariffID)
	if err != nil {
		return nil, fmt.Errorf("could not read tariff id=%d: %w", contract.TariffID, err)
	}
	if trf.Status == tariff.TariffStatusRetired {
		return nil, fmt.Errorf("%w: tariff id=%d is retired", ErrInvalidState, trf.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("could not load contract id=%d: %w", id, err)
	}
	s.log.Info("processing contract", "id", contract.ID)

	if contract.Status != ContractStatusSigned {
		return nil, fmt.Errorf("%w: contract %d is %s", ErrInvalidState, id, contract.Status)
	}
	now := s.clock.Now()
	contract.Status = ContractStatusTerminated
	touchContract(contract, now)
	if err := s.store.UpdateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("could not update contract: %w", err)
	}
	if err := s.audit.Record(ctx, "contract.terminated", contract.ID); err != nil {
		return nil, fmt.Errorf("could not record audit entry: %w", err)
	}

	// TODO(jf): handle partial failures
	s.log.Info("contract terminated", "id", contract.ID)
	return contract, nil
}

type TariffService interface {
	GetTariff(ctx context.Context, id int64) (*tariff.Tariff, error)
}

// ProcessSignedContracts applies terminate to every signed contract and reports how many succeeded.
func (s *Service) ProcessSignedContracts(ctx context.Context) (int, error) {
	pending, err := s.store.ListContractsByStatus(ctx, ContractStatusSigned)
	if err != nil {
		return 0, fmt.Errorf("query signed contracts: %w", err)
	}

	processed := 0
	for _, ctr := range pending {
		_, err := s.TerminateContract(ctx, ctr.ID)
		if err != nil {
			s.log.Error("terminate contract failed", "id", ctr.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}
