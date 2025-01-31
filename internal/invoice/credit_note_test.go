package invoice

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

func (f *fakeStore) FindInvoice(ctx context.Context, id int64) (*Invoice, error) {
	return nil, f.err
}

func (f *fakeStore) ListInvoices(ctx context.Context, limit int, offset int) ([]*Invoice, error) {
	return nil, f.err
}

func (f *fakeStore) ListInvoicesByStatus(ctx context.Context, status InvoiceStatus) ([]*Invoice, error) {
	return nil, f.err
}

func (f *fakeStore) InsertInvoice(ctx context.Context, invoice *Invoice) error {
	return f.err
}

func (f *fakeStore) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	return f.err
}

func (f *fakeStore) GetInvoiceLine(ctx context.Context, id int64) (*InvoiceLine, error) {
	return nil, f.err
}

func (f *fakeStore) ListInvoiceLines(ctx context.Context, limit int, offset int) ([]*InvoiceLine, error) {
	return nil, f.err
}

func (f *fakeStore) ListInvoiceLinesForInvoice(ctx context.Context, invoiceID int64) ([]*InvoiceLine, error) {
	return nil, f.err
}

func (f *fakeStore) InsertInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	return f.err
}

func (f *fakeStore) UpdateInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	return f.err
}

func (f *fakeStore) LoadCreditNote(ctx context.Context, id int64) (*CreditNote, error) {
	return nil, f.err
}

func (f *fakeStore) ListCreditNotes(ctx context.Context, limit int, offset int) ([]*CreditNote, error) {
	return nil, f.err
}

func (f *fakeStore) ListCreditNotesByStatus(ctx context.Context, status CreditNoteStatus) ([]*CreditNote, error) {
	return nil, f.err
}

func (f *fakeStore) ListCreditNotesForInvoice(ctx context.Context, invoiceID int64) ([]*CreditNote, error) {
	return nil, f.err
}

func (f *fakeStore) InsertCreditNote(ctx context.Context, note *CreditNote) error {
	return f.err
}

func (f *fakeStore) UpdateCreditNote(ctx context.Context, note *CreditNote) error {
	return f.err
}

func TestCreateCreditNoteRequest_Validate(t *testing.T) {
	req := CreateCreditNoteRequest{}

	require.Error(t, req.Validate())
}

func TestCreateCreditNoteRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCreditNoteRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateCreditNoteRequest{
				InvoiceID: 1,
				Number:    "Test",
				Reason:    "Test",
				Amount:    1500,
			},
			wantErr: false,
		},
		{
			name:    "missing invoice id",
			req:     CreateCreditNoteRequest{Number: "Test", Reason: "Test", Amount: 1500},
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

func TestService_GetCreditNote_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetCreditNote(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}

func TestCreditNote_CanApply(t *testing.T) {
	note := &CreditNote{}
	require.False(t, note.CanApply())
}
