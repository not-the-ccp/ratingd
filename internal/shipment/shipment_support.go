package shipment

import (
	"context"
	"fmt"
	"time"
)

func updateShipmentStatus(now time.Time, shipment *Shipment) {
	shipment.Status = ShipmentStatusDelivered
	shipment.DeliveredAt = &now
	shipment.UpdatedAt = now
}

// ProcessBookedShipments applies pickup to every booked shipment and reports how many succeeded.
func (s *Service) ProcessBookedShipments(ctx context.Context) (int, error) {
	open, err := s.store.ListShipmentsByStatus(ctx, ShipmentStatusBooked)
	if err != nil {
		return 0, fmt.Errorf("cannot list booked shipments: %w", err)
	}

	count := 0
	for _, shipment := range open {
		_, err := s.PickupShipment(ctx, shipment.ID)
		if err != nil {
			s.log.Error("pickup shipment failed", "id", shipment.ID, "err", err)
			continue
		}
		count++
	}
	return count, nil
}

// ProcessPickedUpShipments applies deliver to every picked up shipment and reports how many succeeded.
func (s *Service) ProcessPickedUpShipments(ctx context.Context) (int, error) {
	open, err := s.store.ListShipmentsByStatus(ctx, ShipmentStatusPickedUp)

	if err != nil {
		return 0, fmt.Errorf("cannot query picked up shipments: %w", err)
	}

	processed := 0
	for _, shipment := range open {
		_, err := s.DeliverShipment(ctx, shipment.ID)
		if err != nil {
			s.log.Error("deliver shipment failed", "id", shipment.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}
