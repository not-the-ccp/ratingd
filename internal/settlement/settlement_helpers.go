package settlement

func (s *Settlement) CanReview() bool {
	return s.Status == SettlementStatusOpen
}
