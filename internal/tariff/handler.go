package tariff

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the tariff API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/tariff", h.getTariff)
	mux.HandleFunc("/api/v1/tariff/create", h.createTariff)
	mux.HandleFunc("/api/v1/tariffs", h.listTariffs)
	mux.HandleFunc("/api/v1/tariff/publish", h.publishTariff)
	mux.HandleFunc("/api/v1/tariff/retire", h.retireTariff)
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

// getTariff handles requests to /api/v1/tariff.
func (h *Handler) getTariff(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	tariff, err := h.svc.GetTariff(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toTariffResponse(tariff))
}

// createTariff handles requests to /api/v1/tariff/create.
func (h *Handler) createTariff(w http.ResponseWriter, r *http.Request) {
	var req CreateTariffRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	tariff, err := h.svc.CreateTariff(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toTariffResponse(tariff))
}

// listTariffs handles requests to /api/v1/tariffs.
func (h *Handler) listTariffs(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)
	tariffs, err := h.svc.ListTariffs(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]TariffResponse, 0, len(tariffs))
	for _, item := range tariffs {
		out = append(out, toTariffResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) publishTariff(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	tariff, err := h.svc.PublishTariff(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toTariffResponse(tariff))
}

// retireTariff handles requests to /api/v1/tariff/retire.
func (h *Handler) retireTariff(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	tariff, err := h.svc.RetireTariff(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toTariffResponse(tariff))
}
