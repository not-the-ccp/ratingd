package carrier

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the carrier API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/carrier", h.getCarrier)
	mux.HandleFunc("/api/v1/carrier/create", h.createCarrier)
	mux.HandleFunc("/api/v1/carriers", h.listCarriers)
	mux.HandleFunc("/api/v1/carrier/approve", h.approveCarrier)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusNotFound, err)
	case errors.Is(err, ErrInvalidState):
		httpx.Error(w, http.StatusConflict, err)
	default:
		httpx.Error(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) getCarrier(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	carrier, err := h.svc.GetCarrier(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCarrierResponse(carrier))
}

func (h *Handler) createCarrier(w http.ResponseWriter, r *http.Request) {
	var req CreateCarrierRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	carrier, err := h.svc.CreateCarrier(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toCarrierResponse(carrier))
}

func (h *Handler) listCarriers(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)
	carriers, err := h.svc.ListCarriers(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]CarrierResponse, 0, len(carriers))
	for _, item := range carriers {
		out = append(out, toCarrierResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) approveCarrier(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	carrier, err := h.svc.ApproveCarrier(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCarrierResponse(carrier))
}
