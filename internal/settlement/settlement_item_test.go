package settlement

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

func (f *fakeStore) FindSettlement(ctx context.Context, id int64) (*Settlement, error) {
	return nil, f.err
}

func (f *fakeStore) ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error) {
	return nil, f.err
}

func (f *fakeStore) ListSettlementsByStatus(ctx context.Context, status SettlementStatus) ([]*Settlement, error) {
	return nil, f.err
}

func (f *fakeStore) InsertSettlement(ctx context.Context, settlement *Settlement) error {
	return f.err
}

func (f *fakeStore) UpdateSettlement(ctx context.Context, settlement *Settlement) error {
	return f.err
}

func (f *fakeStore) GetSettlementItem(ctx context.Context, id int64) (*SettlementItem, error) {
	return nil, f.err
}

func (f *fakeStore) ListSettlementItems(ctx context.Context, limit int, offset int) ([]*SettlementItem, error) {
	return nil, f.err
}

func (f *fakeStore) ListSettlementItemsForSettlement(ctx context.Context, settlementID int64) ([]*SettlementItem, error) {
	return nil, f.err
}

func (f *fakeStore) InsertSettlementItem(ctx context.Context, item *SettlementItem) error {
	return f.err
}

func (f *fakeStore) UpdateSettlementItem(ctx context.Context, item *SettlementItem) error {
	return f.err
}

func TestCreateSettlementItemRequest_Validate(t *testing.T) {
	req := CreateSettlementItemRequest{}

	require.Error(t, req.Validate())
}

func TestCreateSettlementItemRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateSettlementItemRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateSettlementItemRequest{
				SettlementID: 1,
				ShipmentID:   1,
				Amount:       1,
				Note:         "Sample",
			},
			wantErr: false,
		},
		{
			name:    "missing shipment id",
			req:     CreateSettlementItemRequest{SettlementID: 1, Amount: 1, Note: "Sample"},
			wantErr: true,
		},
		{
			name:    "missing note",
			req:     CreateSettlementItemRequest{SettlementID: 1, ShipmentID: 1, Amount: 1},
			wantErr: true,
		},
		{
			name:    "missing settlement id",
			req:     CreateSettlementItemRequest{ShipmentID: 1, Amount: 1, Note: "Sample"},
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

func TestService_GetSettlementItem_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetSettlementItem(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}

func TestSettlementItem_HasAmount(t *testing.T) {
	item := &SettlementItem{}
	require.False(t, item.HasAmount())
}
