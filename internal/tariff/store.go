package tariff

import (
	"context"
)

// Store persists tariff data.
type Store interface {
	FindTariff(ctx context.Context, id int64) (*Tariff, error)
	ListTariffs(ctx context.Context, limit int, offset int) ([]*Tariff, error)
	ListTariffsByStatus(ctx context.Context, status TariffStatus) ([]*Tariff, error)
	InsertTariff(ctx context.Context, tariff *Tariff) error
	UpdateTariff(ctx context.Context, tariff *Tariff) error
	FindWeightBreak(ctx context.Context, id int64) (*WeightBreak, error)
	ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error)
	ListWeightBreaksForTariff(ctx context.Context, tariffID int64) ([]*WeightBreak, error)
	InsertWeightBreak(ctx context.Context, breakVal *WeightBreak) error
	UpdateWeightBreak(ctx context.Context, breakVal *WeightBreak) error
}
