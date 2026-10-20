package contract

import (
	"context"
)

// Store persists contract data.
type Store interface {
	FetchContract(ctx context.Context, id int64) (*Contract, error)
	ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error)
	ListContractsByStatus(ctx context.Context, status ContractStatus) ([]*Contract, error)
	InsertContract(ctx context.Context, contract *Contract) error
	UpdateContract(ctx context.Context, contract *Contract) error
}
