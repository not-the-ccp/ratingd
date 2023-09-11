package carrier

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getCarrier_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/carrier?id=abc", nil)

	h.getCarrier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
