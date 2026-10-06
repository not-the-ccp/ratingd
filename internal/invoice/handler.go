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
	mux.HandleFunc("/api/v1/invoices", h.listInvoices)
	mux.HandleFunc("/api/v1/invoice/issue", h.issueInvoice)
	mux.HandleFunc("/api/v1/invoice/mark-paid", h.markPaidInvoice)
	mux.HandleFunc("/api/v1/invoice/void", h.voidInvoice)
	mux.HandleFunc("/api/v1/invoice-line", h.getInvoiceLine)
	mux.HandleFunc("/api/v1/invoice-line/create", h.createInvoiceLine)
	mux.HandleFunc("/api/v1/invoice-lines", h.listInvoiceLines)
	mux.HandleFunc("/api/v1/credit-note", h.getCreditNote)
	mux.HandleFunc("/api/v1/credit-note/create", h.createCreditNote)
	mux.HandleFunc("/api/v1/credit-note/apply", h.applyCreditNote)
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

// getInvoiceLine handles requests to /api/v1/invoice-line.
func (h *Handler) getInvoiceLine(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	line, err := h.svc.GetInvoiceLine(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toInvoiceLineResponse(line))
}

func (h *Handler) createInvoiceLine(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceLineRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	line, err := h.svc.CreateInvoiceLine(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toInvoiceLineResponse(line))
}

func (h *Handler) listInvoiceLines(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	lines, err := h.svc.ListInvoiceLines(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]InvoiceLineResponse, 0, len(lines))
	for _, item := range lines {
		out = append(out, toInvoiceLineResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) getCreditNote(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	note, err := h.svc.GetCreditNote(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCreditNoteResponse(note))
}

// createCreditNote handles requests to /api/v1/credit-note/create.
func (h *Handler) createCreditNote(w http.ResponseWriter, r *http.Request) {
	var req CreateCreditNoteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	note, err := h.svc.CreateCreditNote(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toCreditNoteResponse(note))
}

// applyCreditNote handles requests to /api/v1/credit-note/apply.
func (h *Handler) applyCreditNote(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	note, err := h.svc.ApplyCreditNote(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCreditNoteResponse(note))
}
