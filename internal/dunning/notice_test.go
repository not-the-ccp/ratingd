package dunning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotice_CanSend(t *testing.T) {
	notice := &Notice{}
	require.False(t, notice.CanSend())
}
