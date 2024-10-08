package rating

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuote_CanAccept(t *testing.T) {
	quote := &Quote{}
	require.False(t, quote.CanAccept())
}
