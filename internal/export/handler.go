package export

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the export API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/batch", h.getBatch)
	mux.HandleFunc("/api/v1/batch/create", h.createBatch)
	mux.HandleFunc("/api/v1/batches", h.listBatches)
	mux.HandleFunc("/api/v1/batch/finish", h.finishBatch)
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

func (h *Handler) getBatch(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	batch, err := h.svc.GetBatch(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toBatchResponse(batch))
}

func (h *Handler) createBatch(w http.ResponseWriter, r *http.Request) {
	var req CreateBatchRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	batch, err := h.svc.CreateBatch(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toBatchResponse(batch))
}

func (h *Handler) listBatches(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	batches, err := h.svc.ListBatches(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]BatchResponse, 0, len(batches))
	for _, item := range batches {
		out = append(out, toBatchResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

// finishBatch handles requests to /api/v1/batch/finish.
func (h *Handler) finishBatch(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	batch, err := h.svc.FinishBatch(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toBatchResponse(batch))
}
