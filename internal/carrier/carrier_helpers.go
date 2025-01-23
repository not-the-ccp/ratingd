package carrier

func (c *Carrier) CanApprove() bool {
	return c.Status == CarrierStatusOnboarding
}
