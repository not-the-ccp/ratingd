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
	mux.HandleFunc("/api/v1/tariff/retire", h.retireTariff)
	mux.HandleFunc("/api/v1/weight-break", h.getWeightBreak)
	mux.HandleFunc("/api/v1/weight-break/create", h.createWeightBreak)
	mux.HandleFunc("/api/v1/weight-breaks", h.listWeightBreaks)
	mux.HandleFunc("/api/v1/fuel-surcharge", h.getFuelSurcharge)
	mux.HandleFunc("/api/v1/fuel-surcharge/create", h.createFuelSurcharge)
	mux.HandleFunc("/api/v1/fuel-surcharges", h.listFuelSurcharges)
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

func (h *Handler) getWeightBreak(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	breakEntry, err := h.svc.GetWeightBreak(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, toWeightBreakResponse(breakEntry))
}

func (h *Handler) createWeightBreak(w http.ResponseWriter, r *http.Request) {
	var req CreateWeightBreakRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	weightBreak, err := h.svc.CreateWeightBreak(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toWeightBreakResponse(weightBreak))
}

func (h *Handler) listWeightBreaks(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)
	weightBreaks, err := h.svc.ListWeightBreaks(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]WeightBreakResponse, 0, len(weightBreaks))
	for _, item := range weightBreaks {
		out = append(out, toWeightBreakResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) getFuelSurcharge(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	surcharge, err := h.svc.GetFuelSurcharge(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, toFuelSurchargeResponse(surcharge))
}

func (h *Handler) createFuelSurcharge(w http.ResponseWriter, r *http.Request) {
	var req CreateFuelSurchargeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	surcharge, err := h.svc.CreateFuelSurcharge(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toFuelSurchargeResponse(surcharge))
}

// listFuelSurcharges handles requests to /api/v1/fuel-surcharges.
func (h *Handler) listFuelSurcharges(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)
	surcharges, err := h.svc.ListFuelSurcharges(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]FuelSurchargeResponse, 0, len(surcharges))
	for _, item := range surcharges {
		out = append(out, toFuelSurchargeResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}
