package rating

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getSurcharge_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/surcharge?id=abc", nil)

	h.getSurcharge(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
