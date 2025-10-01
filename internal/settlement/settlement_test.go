package settlement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSettlement_CanRelease(t *testing.T) {
	settlement := &Settlement{}
	require.False(t, settlement.CanRelease())
}

func TestSettlement_CanReview(t *testing.T) {
	settlement := &Settlement{}
	require.False(t, settlement.CanReview())
}

func TestSettlement_HasDeductions(t *testing.T) {
	settlement := &Settlement{}
	require.False(t, settlement.HasDeductions())
}

func TestSettlement_HasGrossAmount(t *testing.T) {
	settlement := &Settlement{}
	require.False(t, settlement.HasGrossAmount())
}
