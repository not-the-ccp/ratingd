package settlement

import (
	"context"
)

// Store persists settlement data.
type Store interface {
	FindSettlement(ctx context.Context, id int64) (*Settlement, error)
	ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error)
	ListSettlementsByStatus(ctx context.Context, status SettlementStatus) ([]*Settlement, error)
	InsertSettlement(ctx context.Context, settlement *Settlement) error
	UpdateSettlement(ctx context.Context, settlement *Settlement) error
}
