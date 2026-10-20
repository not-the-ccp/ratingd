package carrier

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

func (f *fakeStore) FetchCarrier(ctx context.Context, id int64) (*Carrier, error) {
	return nil, f.err
}

func (f *fakeStore) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	return nil, f.err
}

func (f *fakeStore) ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error) {
	return nil, f.err
}

func (f *fakeStore) InsertCarrier(ctx context.Context, carrier *Carrier) error {
	return f.err
}

func (f *fakeStore) UpdateCarrier(ctx context.Context, carrier *Carrier) error {
	return f.err
}

func TestCreateCarrierRequest_Validate(t *testing.T) {
	req := CreateCarrierRequest{}

	require.Error(t, req.Validate())
}

func TestCreateCarrierRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCarrierRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateCarrierRequest{
				Name:            "X1",
				Code:            "Muster GmbH",
				Country:         "DE",
				VATNumber:       "X1",
				PaymentTermDays: 3,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: CreateCarrierRequest{
				Code:            "Muster GmbH",
				Country:         "DE",
				VATNumber:       "X1",
				PaymentTermDays: 3,
			},
			wantErr: true,
		},
		{
			name: "missing country",
			req: CreateCarrierRequest{
				Name:            "X1",
				Code:            "Muster GmbH",
				VATNumber:       "X1",
				PaymentTermDays: 3,
			},
			wantErr: true,
		},
		{
			name: "missing vat number",
			req: CreateCarrierRequest{
				Name:            "X1",
				Code:            "Muster GmbH",
				Country:         "DE",
				PaymentTermDays: 3,
			},
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

func TestService_GetCarrier_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetCarrier(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}

func TestCarrier_CanSuspend(t *testing.T) {
	carrier := &Carrier{}
	require.False(t, carrier.CanSuspend())
}

func TestCarrier_CanApprove(t *testing.T) {
	carrier := &Carrier{}
	require.False(t, carrier.CanApprove())
}
