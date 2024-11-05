package invoice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoice_CanIssue(t *testing.T) {
	invoice := &Invoice{}
	require.False(t, invoice.CanIssue())
}
