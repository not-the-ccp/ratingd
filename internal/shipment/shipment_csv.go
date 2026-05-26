package shipment

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeShipmentsCSV(w io.Writer, items []*Shipment) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{
		"id",
		"origin_zip",
		"pallets",
		"pickup_at",
		"origin_country",
	}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, shipment := range items {
		rec := []string{
			strconv.FormatInt(shipment.ID, 10),
			shipment.OriginZip,
			strconv.Itoa(shipment.Pallets),
			shipment.PickupAt.Format(time.RFC3339),
			shipment.OriginCountry,
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("write shipment %d: %w", shipment.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

func (s *Service) ExportShipments(ctx context.Context, w io.Writer, limit int) error {
	shipments, err := s.store.ListShipments(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("list shipments: %w", err)
	}

	return writeShipmentsCSV(w, shipments)
}
