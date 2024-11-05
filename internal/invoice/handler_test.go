package invoice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getCreditNote_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/credit-note?id=abc", nil)

	h.getCreditNote(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
