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
}
