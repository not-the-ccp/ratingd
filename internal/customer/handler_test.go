package customer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getBillingAddress_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/billing-address?id=abc", nil)

	h.getBillingAddress(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
