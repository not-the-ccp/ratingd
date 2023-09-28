package tariff

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getTariff_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tariff?id=abc", nil)

	h.getTariff(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
