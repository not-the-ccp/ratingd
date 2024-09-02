package rating

import (
	"context"
)

// Store persists rating data.
type Store interface {
	GetQuote(ctx context.Context, id int64) (*Quote, error)
	ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error)
	ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error)
	InsertQuote(ctx context.Context, quote *Quote) error
	UpdateQuote(ctx context.Context, quote *Quote) error
}
