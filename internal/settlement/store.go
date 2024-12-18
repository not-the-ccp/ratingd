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
	GetSettlementItem(ctx context.Context, id int64) (*SettlementItem, error)
	ListSettlementItems(ctx context.Context, limit int, offset int) ([]*SettlementItem, error)
	ListSettlementItemsForSettlement(ctx context.Context, settlementID int64) ([]*SettlementItem, error)
	InsertSettlementItem(ctx context.Context, item *SettlementItem) error
	UpdateSettlementItem(ctx context.Context, item *SettlementItem) error
}
