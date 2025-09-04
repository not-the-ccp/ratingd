package invoice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoice_CanIssue(t *testing.T) {
	invoice := &Invoice{}
	require.False(t, invoice.CanIssue())
}

func TestInvoice_CanMarkPaid(t *testing.T) {
	invoice := &Invoice{}
	require.False(t, invoice.CanMarkPaid())
}

func TestInvoice_HasTaxAmount(t *testing.T) {
	invoice := &Invoice{}
	require.False(t, invoice.HasTaxAmount())
}

func TestInvoice_HasNetAmount(t *testing.T) {
	invoice := &Invoice{}
	require.False(t, invoice.HasNetAmount())
}
