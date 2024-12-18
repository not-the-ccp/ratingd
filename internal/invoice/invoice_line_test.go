package invoice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoiceLine_HasUnitPrice(t *testing.T) {
	line := &InvoiceLine{}
	require.False(t, line.HasUnitPrice())
}
