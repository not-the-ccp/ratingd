package rating

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the rating API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/quote", h.getQuote)
	mux.HandleFunc("/api/v1/quote/create", h.createQuote)
	mux.HandleFunc("/api/v1/quotes", h.listQuotes)
	mux.HandleFunc("/api/v1/quote/accept", h.acceptQuote)
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

func (h *Handler) getQuote(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	quote, err := h.svc.GetQuote(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toQuoteResponse(quote))
}

// createQuote handles requests to /api/v1/quote/create.
func (h *Handler) createQuote(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	quote, err := h.svc.CreateQuote(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toQuoteResponse(quote))
}

// listQuotes handles requests to /api/v1/quotes.
func (h *Handler) listQuotes(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	quotes, err := h.svc.ListQuotes(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]QuoteResponse, 0, len(quotes))
	for _, item := range quotes {
		out = append(out, toQuoteResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) acceptQuote(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	quote, err := h.svc.AcceptQuote(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toQuoteResponse(quote))
}
