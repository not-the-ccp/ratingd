package shipment

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

func (f *fakeStore) FetchShipment(ctx context.Context, id int64) (*Shipment, error) {
	return nil, f.err
}

func (f *fakeStore) ListShipments(ctx context.Context, limit int, offset int) ([]*Shipment, error) {
	return nil, f.err
}

func (f *fakeStore) ListShipmentsByStatus(ctx context.Context, status ShipmentStatus) ([]*Shipment, error) {
	return nil, f.err
}

func (f *fakeStore) InsertShipment(ctx context.Context, shipment *Shipment) error {
	return f.err
}

func (f *fakeStore) UpdateShipment(ctx context.Context, shipment *Shipment) error {
	return f.err
}

func (f *fakeStore) FetchTrackingEvent(ctx context.Context, id int64) (*TrackingEvent, error) {
	return nil, f.err
}

func (f *fakeStore) ListTrackingEvents(ctx context.Context, limit int, offset int) ([]*TrackingEvent, error) {
	return nil, f.err
}

func (f *fakeStore) ListTrackingEventsForShipment(ctx context.Context, shipmentID int64) ([]*TrackingEvent, error) {
	return nil, f.err
}

func (f *fakeStore) InsertTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	return f.err
}

func (f *fakeStore) UpdateTrackingEvent(ctx context.Context, event *TrackingEvent) error {
	return f.err
}

func TestCreateTrackingEventRequest_Validate(t *testing.T) {
	req := CreateTrackingEventRequest{}

	require.Error(t, req.Validate())
}

func TestCreateTrackingEventRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTrackingEventRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     CreateTrackingEventRequest{ShipmentID: 1, Code: "ABC-1", Location: "X1"},
			wantErr: false,
		},
		{
			name:    "missing shipment id",
			req:     CreateTrackingEventRequest{Code: "ABC-1", Location: "X1"},
			wantErr: true,
		},
		{
			name:    "missing code",
			req:     CreateTrackingEventRequest{ShipmentID: 1, Location: "X1"},
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

func TestService_GetTrackingEvent_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetTrackingEvent(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}
