package customer

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the customer API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/customer", h.getCustomer)
	mux.HandleFunc("/api/v1/customers", h.listCustomers)
	mux.HandleFunc("/api/v1/customer/activate", h.activateCustomer)
	mux.HandleFunc("/api/v1/billing-address", h.getBillingAddress)
	mux.HandleFunc("/api/v1/billing-addresses", h.listBillingAddresses)
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

func (h *Handler) getCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	customer, err := h.svc.GetCustomer(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCustomerResponse(customer))
}

func (h *Handler) listCustomers(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)
	customers, err := h.svc.ListCustomers(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]CustomerResponse, 0, len(customers))
	for _, item := range customers {
		out = append(out, toCustomerResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) activateCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	customer, err := h.svc.ActivateCustomer(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCustomerResponse(customer))
}

// getBillingAddress handles requests to /api/v1/billing-address.
func (h *Handler) getBillingAddress(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	address, err := h.svc.GetBillingAddress(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, toBillingAddressResponse(address))
}

func (h *Handler) listBillingAddresses(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	addresses, err := h.svc.ListBillingAddresses(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]BillingAddressResponse, 0, len(addresses))
	for _, item := range addresses {
		out = append(out, toBillingAddressResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}
