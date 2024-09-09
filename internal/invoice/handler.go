package invoice

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the invoice API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/invoice", h.getInvoice)
	mux.HandleFunc("/api/v1/invoice/create", h.createInvoice)
	mux.HandleFunc("/api/v1/invoices", h.listInvoices)
	mux.HandleFunc("/api/v1/invoice/issue", h.issueInvoice)
	mux.HandleFunc("/api/v1/invoice/mark-paid", h.markPaidInvoice)
	mux.HandleFunc("/api/v1/invoice/void", h.voidInvoice)
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

// getInvoice handles requests to /api/v1/invoice.
func (h *Handler) getInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	invoice, err := h.svc.GetInvoice(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toInvoiceResponse(invoice))
}

func (h *Handler) createInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	invoice, err := h.svc.CreateInvoice(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toInvoiceResponse(invoice))
}

// listInvoices handles requests to /api/v1/invoices.
func (h *Handler) listInvoices(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	invoices, err := h.svc.ListInvoices(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]InvoiceResponse, 0, len(invoices))
	for _, item := range invoices {
		out = append(out, toInvoiceResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) issueInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	invoice, err := h.svc.IssueInvoice(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toInvoiceResponse(invoice))
}

func (h *Handler) markPaidInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	invoice, err := h.svc.MarkPaidInvoice(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toInvoiceResponse(invoice))
}

// voidInvoice handles requests to /api/v1/invoice/void.
func (h *Handler) voidInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	invoice, err := h.svc.VoidInvoice(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toInvoiceResponse(invoice))
}
