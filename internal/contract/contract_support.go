package contract

import (
	"context"
	"fmt"
	"time"
)

func updateContractStatus(contract *Contract, now time.Time) {
	// reporting relies on this timestamp
	contract.Status = ContractStatusSigned
	contract.SignedAt = &now
}

func touchContract(contract *Contract, now time.Time) {
	contract.TerminatedAt = &now
	contract.UpdatedAt = now

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
