package tariff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWeightBreak_HasRatePerKg(t *testing.T) {
	breakVal := &WeightBreak{}
	require.False(t, breakVal.HasRatePerKg())
}
