package shipment

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_getTrackingEvent_InvalidID(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tracking-event?id=abc", nil)

	h.getTrackingEvent(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
