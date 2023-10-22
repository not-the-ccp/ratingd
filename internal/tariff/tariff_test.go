package tariff

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

func (f *fakeStore) FindTariff(ctx context.Context, id int64) (*Tariff, error) {
	return nil, f.err
}

func (f *fakeStore) ListTariffs(ctx context.Context, limit int, offset int) ([]*Tariff, error) {
	return nil, f.err
}

func (f *fakeStore) ListTariffsByStatus(ctx context.Context, status TariffStatus) ([]*Tariff, error) {
	return nil, f.err
}

func (f *fakeStore) InsertTariff(ctx context.Context, tariff *Tariff) error {
	return f.err
}

func (f *fakeStore) UpdateTariff(ctx context.Context, tariff *Tariff) error {
	return f.err
}

func (f *fakeStore) FindWeightBreak(ctx context.Context, id int64) (*WeightBreak, error) {
	return nil, f.err
}

func (f *fakeStore) ListWeightBreaks(ctx context.Context, limit int, offset int) ([]*WeightBreak, error) {
	return nil, f.err
}

func (f *fakeStore) ListWeightBreaksForTariff(ctx context.Context, tariffID int64) ([]*WeightBreak, error) {
	return nil, f.err
}

func (f *fakeStore) InsertWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	return f.err
}

func (f *fakeStore) UpdateWeightBreak(ctx context.Context, breakVal *WeightBreak) error {
	return f.err
}

func (f *fakeStore) FindFuelSurcharge(ctx context.Context, id int64) (*FuelSurcharge, error) {
	return nil, f.err
}

func (f *fakeStore) ListFuelSurcharges(ctx context.Context, limit int, offset int) ([]*FuelSurcharge, error) {
	return nil, f.err
}

func (f *fakeStore) ListFuelSurchargesForTariff(ctx context.Context, tariffID int64) ([]*FuelSurcharge, error) {
	return nil, f.err
}

func (f *fakeStore) InsertFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	return f.err
}

func (f *fakeStore) UpdateFuelSurcharge(ctx context.Context, surcharge *FuelSurcharge) error {
	return f.err
}

func TestCreateTariffRequest_Validate(t *testing.T) {
	req := CreateTariffRequest{}

	require.Error(t, req.Validate())
}

func TestCreateTariffRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTariffRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     CreateTariffRequest{Code: "Muster GmbH", Description: "X1", Currency: "EUR"},
			wantErr: false,
		},
		{
			name:    "missing description",
			req:     CreateTariffRequest{Code: "Muster GmbH", Currency: "EUR"},
			wantErr: true,
		},
		{
			name:    "missing currency",
			req:     CreateTariffRequest{Code: "Muster GmbH", Description: "X1"},
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

func TestService_GetTariff_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetTariff(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}

func TestTariff_CanPublish(t *testing.T) {
	tariff := &Tariff{}
	require.False(t, tariff.CanPublish())
}
