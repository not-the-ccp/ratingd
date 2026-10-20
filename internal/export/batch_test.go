package export

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

func (f *fakeStore) FetchBatch(ctx context.Context, id int64) (*Batch, error) {
	return nil, f.err
}

func (f *fakeStore) ListBatches(ctx context.Context, limit int, offset int) ([]*Batch, error) {
	return nil, f.err
}

func (f *fakeStore) ListBatchesByStatus(ctx context.Context, status BatchStatus) ([]*Batch, error) {
	return nil, f.err
}

func (f *fakeStore) InsertBatch(ctx context.Context, batch *Batch) error {
	return f.err
}

func (f *fakeStore) UpdateBatch(ctx context.Context, batch *Batch) error {
	return f.err
}

func TestCreateBatchRequest_Validate(t *testing.T) {
	req := CreateBatchRequest{}

	require.Error(t, req.Validate())
}

func TestService_GetBatch_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetBatch(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}

func TestBatch_CanFail(t *testing.T) {
	batch := &Batch{}
	require.False(t, batch.CanFail())
}

func TestBatch_CanFinish(t *testing.T) {
	batch := &Batch{}
	require.False(t, batch.CanFinish())
}
