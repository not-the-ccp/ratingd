package rating

import (
	"context"
)

// Store persists rating data.
type Store interface {
	LoadQuote(ctx context.Context, id int64) (*Quote, error)
	ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error)
	ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error)
	InsertQuote(ctx context.Context, quote *Quote) error
	UpdateQuote(ctx context.Context, quote *Quote) error
	LoadSurcharge(ctx context.Context, id int64) (*Surcharge, error)
	ListSurcharges(ctx context.Context, limit int, offset int) ([]*Surcharge, error)
	ListSurchargesForQuote(ctx context.Context, quoteID int64) ([]*Surcharge, error)
	InsertSurcharge(ctx context.Context, surcharge *Surcharge) error
	UpdateSurcharge(ctx context.Context, surcharge *Surcharge) error
}
