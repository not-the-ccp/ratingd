package settlement

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the settlement API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/settlement", h.getSettlement)
	mux.HandleFunc("/api/v1/settlement/create", h.createSettlement)
	mux.HandleFunc("/api/v1/settlements", h.listSettlements)
	mux.HandleFunc("/api/v1/settlement/review", h.reviewSettlement)
	mux.HandleFunc("/api/v1/settlement-item", h.getSettlementItem)
	mux.HandleFunc("/api/v1/settlement-item/create", h.createSettlementItem)
	mux.HandleFunc("/api/v1/settlement-items", h.listSettlementItems)
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

// getSettlement handles requests to /api/v1/settlement.
func (h *Handler) getSettlement(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	settlement, err := h.svc.GetSettlement(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, toSettlementResponse(settlement))
}

func (h *Handler) createSettlement(w http.ResponseWriter, r *http.Request) {
	var req CreateSettlementRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	settlement, err := h.svc.CreateSettlement(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toSettlementResponse(settlement))
}

func (h *Handler) listSettlements(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	settlements, err := h.svc.ListSettlements(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]SettlementResponse, 0, len(settlements))
	for _, item := range settlements {
		out = append(out, toSettlementResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) reviewSettlement(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	settlement, err := h.svc.ReviewSettlement(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toSettlementResponse(settlement))
}

func (h *Handler) getSettlementItem(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.GetSettlementItem(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, toSettlementItemResponse(item))
}

func (h *Handler) createSettlementItem(w http.ResponseWriter, r *http.Request) {
	var req CreateSettlementItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.CreateSettlementItem(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toSettlementItemResponse(item))
}

// listSettlementItems handles requests to /api/v1/settlement-items.
func (h *Handler) listSettlementItems(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	settlementItems, err := h.svc.ListSettlementItems(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]SettlementItemResponse, 0, len(settlementItems))
	for _, item := range settlementItems {
		out = append(out, toSettlementItemResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}
