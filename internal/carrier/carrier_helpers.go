package carrier

import (
	"fmt"
	"strings"
)

func (c *Carrier) CanApprove() bool {
	return c.Status == CarrierStatusOnboarding
}

func (c *Carrier) CanSuspend() bool {
	return c.Status == CarrierStatusApproved
}

// ParseCarrierStatus converts a raw status string, ignoring case and surrounding space.
func ParseCarrierStatus(s string) (CarrierStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "onboarding":
		return CarrierStatusOnboarding, nil
	case "approved":
		return CarrierStatusApproved, nil
	case "suspended":
		return CarrierStatusSuspended, nil
	default:
		return "", fmt.Errorf("unknown carrier status %q", s)
	}
}
