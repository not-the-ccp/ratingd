package invoice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoiceLine_HasUnitPrice(t *testing.T) {
	line := &InvoiceLine{}
	require.False(t, line.HasUnitPrice())
}

func TestInvoiceLine_HasAmount(t *testing.T) {
	line := &InvoiceLine{}
	require.False(t, line.HasAmount())
}

func TestInvoiceLine_HasDiscountAmount(t *testing.T) {
	line := &InvoiceLine{}
	require.False(t, line.HasDiscountAmount())
}
