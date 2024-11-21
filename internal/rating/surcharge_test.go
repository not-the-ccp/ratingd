package rating

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeStore implements Store and fails every call with err.
type fakeStore struct {
	err error
}

func (f *fakeStore) GetQuote(ctx context.Context, id int64) (*Quote, error) {
	return nil, f.err
}

func (f *fakeStore) ListQuotes(ctx context.Context, limit int, offset int) ([]*Quote, error) {
	return nil, f.err
}

func (f *fakeStore) ListQuotesByStatus(ctx context.Context, status QuoteStatus) ([]*Quote, error) {
	return nil, f.err
}

func (f *fakeStore) InsertQuote(ctx context.Context, quote *Quote) error {
	return f.err
}

func (f *fakeStore) UpdateQuote(ctx context.Context, quote *Quote) error {
	return f.err
}

func (f *fakeStore) FindSurcharge(ctx context.Context, id int64) (*Surcharge, error) {
	return nil, f.err
}

func (f *fakeStore) ListSurcharges(ctx context.Context, limit int, offset int) ([]*Surcharge, error) {
	return nil, f.err
}

func (f *fakeStore) ListSurchargesForQuote(ctx context.Context, quoteID int64) ([]*Surcharge, error) {
	return nil, f.err
}

func (f *fakeStore) InsertSurcharge(ctx context.Context, surcharge *Surcharge) error {
	return f.err
}

func (f *fakeStore) UpdateSurcharge(ctx context.Context, surcharge *Surcharge) error {
	return f.err
}

func TestCreateSurchargeRequest_Validate(t *testing.T) {
	req := CreateSurchargeRequest{}

	require.Error(t, req.Validate())
}

func TestCreateSurchargeRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateSurchargeRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateSurchargeRequest{
				QuoteID:     1,
				Code:        "Muster GmbH",
				Description: "Sample",
				Amount:      1500,
			},
			wantErr: false,
		},
		{
			name:    "missing description",
			req:     CreateSurchargeRequest{QuoteID: 1, Code: "Muster GmbH", Amount: 1500},
			wantErr: true,
		},
		{
			name:    "missing code",
			req:     CreateSurchargeRequest{QuoteID: 1, Description: "Sample", Amount: 1500},
			wantErr: true,
		},
		{
			name:    "missing quote id",
			req:     CreateSurchargeRequest{Code: "Muster GmbH", Description: "Sample", Amount: 1500},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestService_GetSurcharge_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetSurcharge(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}

func TestSurcharge_HasAmount(t *testing.T) {
	surcharge := &Surcharge{}
	require.False(t, surcharge.HasAmount())
}
