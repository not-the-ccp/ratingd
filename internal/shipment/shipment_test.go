package shipment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShipment_CanPickup(t *testing.T) {
	shipment := &Shipment{}
	require.False(t, shipment.CanPickup())
}

func TestShipment_CanDeliver(t *testing.T) {
	shipment := &Shipment{}
	require.False(t, shipment.CanDeliver())
}
