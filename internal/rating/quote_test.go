package rating

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuote_CanAccept(t *testing.T) {
	quote := &Quote{}
	require.False(t, quote.CanAccept())
}

func TestQuote_HasNetAmount(t *testing.T) {
	quote := &Quote{}
	require.False(t, quote.HasNetAmount())
}

func TestQuote_CanSupersede(t *testing.T) {
	quote := &Quote{}
	require.False(t, quote.CanSupersede())
}

func TestQuote_HasFuelAmount(t *testing.T) {
	quote := &Quote{}
	require.False(t, quote.HasFuelAmount())
}

func TestQuote_HasDiscountAmount(t *testing.T) {
	quote := &Quote{}
	require.False(t, quote.HasDiscountAmount())
}
