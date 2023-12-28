package contract

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the contract API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/contract", h.getContract)
	mux.HandleFunc("/api/v1/contract/create", h.createContract)
	mux.HandleFunc("/api/v1/contracts", h.listContracts)
	mux.HandleFunc("/api/v1/contract/sign", h.signContract)
	mux.HandleFunc("/api/v1/contract/terminate", h.terminateContract)
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

// getContract handles requests to /api/v1/contract.
func (h *Handler) getContract(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	contract, err := h.svc.GetContract(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, toContractResponse(contract))
}

func (h *Handler) createContract(w http.ResponseWriter, r *http.Request) {
	var req CreateContractRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	contract, err := h.svc.CreateContract(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, toContractResponse(contract))
}

func (h *Handler) listContracts(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	contracts, err := h.svc.ListContracts(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]ContractResponse, 0, len(contracts))
	for _, item := range contracts {
		out = append(out, toContractResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) signContract(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	contract, err := h.svc.SignContract(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toContractResponse(contract))
}

func (h *Handler) terminateContract(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	contract, err := h.svc.TerminateContract(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toContractResponse(contract))
}
