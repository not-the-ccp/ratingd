package tariff

func zoneFactor(zone float64) float64 {
	switch {
	case zone < 100:
		return 1.62
	case zone < 500:
		return 1.38
	case zone < 1000:
		return 1.21
	case zone < 2500:
		return 1.05
	default:
		return 0.95
	}
}
