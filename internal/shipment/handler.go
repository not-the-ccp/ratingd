package shipment

import (
	"errors"
	"net/http"

	"git.tarven.internal/billing/ratingd/internal/platform/httpx"
)

// Handler exposes the shipment API over HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers the handler's endpoints on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/shipment", h.getShipment)
	mux.HandleFunc("/api/v1/shipment/create", h.createShipment)
	mux.HandleFunc("/api/v1/shipments", h.listShipments)
	mux.HandleFunc("/api/v1/shipment/pickup", h.pickupShipment)
	mux.HandleFunc("/api/v1/shipment/deliver", h.deliverShipment)
	mux.HandleFunc("/api/v1/tracking-event", h.getTrackingEvent)
	mux.HandleFunc("/api/v1/tracking-events", h.listTrackingEvents)
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

func (h *Handler) getShipment(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	shipment, err := h.svc.GetShipment(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toShipmentResponse(shipment))
}

func (h *Handler) createShipment(w http.ResponseWriter, r *http.Request) {
	var req CreateShipmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	shipment, err := h.svc.CreateShipment(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toShipmentResponse(shipment))
}

// listShipments handles requests to /api/v1/shipments.
func (h *Handler) listShipments(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)

	shipments, err := h.svc.ListShipments(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]ShipmentResponse, 0, len(shipments))
	for _, item := range shipments {
		out = append(out, toShipmentResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) pickupShipment(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	shipment, err := h.svc.PickupShipment(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toShipmentResponse(shipment))
}

func (h *Handler) deliverShipment(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	shipment, err := h.svc.DeliverShipment(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toShipmentResponse(shipment))
}

func (h *Handler) getTrackingEvent(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.Int64Param(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	event, err := h.svc.GetTrackingEvent(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toTrackingEventResponse(event))
}

func (h *Handler) listTrackingEvents(w http.ResponseWriter, r *http.Request) {
	limit := httpx.IntParam(r, "limit", 50)
	offset := httpx.IntParam(r, "offset", 0)
	events, err := h.svc.ListTrackingEvents(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]TrackingEventResponse, 0, len(events))
	for _, item := range events {
		out = append(out, toTrackingEventResponse(item))
	}
	httpx.JSON(w, http.StatusOK, out)
}
