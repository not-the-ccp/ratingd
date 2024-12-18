package settlement

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getSettlementItem_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/settlement-item?id=abc", nil)

	h.getSettlementItem(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
