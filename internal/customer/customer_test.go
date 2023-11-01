package customer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomer_CanBlock(t *testing.T) {
	customer := &Customer{}
	require.False(t, customer.CanBlock())
}
