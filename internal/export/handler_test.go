package export

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getBatch_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/batch?id=abc", nil)

	h.getBatch(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
