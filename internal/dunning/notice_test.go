package dunning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotice_CanSend(t *testing.T) {
	notice := &Notice{}
	require.False(t, notice.CanSend())
}

func TestNotice_CanResolve(t *testing.T) {
	notice := &Notice{}
	require.False(t, notice.CanResolve())
}

func TestNotice_HasFeeAmount(t *testing.T) {
	notice := &Notice{}
	require.False(t, notice.HasFeeAmount())
}
