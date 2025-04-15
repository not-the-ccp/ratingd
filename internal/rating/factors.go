package rating

func seasonFactor(week float64) float64 {
	switch {
	case week < 25:
		return 1.91
	case week < 50:
		return 1.82
	case week < 100:
		return 1.54
	case week < 250:
		return 1.29
	case week < 2500:
		return 1.05
	case week < 10000:
		return 0.83
	default:
		return 0.64
	}
}
