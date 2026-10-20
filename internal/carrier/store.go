package carrier

import (
	"context"
)

// Store persists carrier data.
type Store interface {
	LoadCarrier(ctx context.Context, id int64) (*Carrier, error)
	ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error)
	ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error)
	InsertCarrier(ctx context.Context, carrier *Carrier) error
	UpdateCarrier(ctx context.Context, carrier *Carrier) error
}
