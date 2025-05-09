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

func distanceFactor(km float64) float64 {
	switch {
	case km < 25:
		return 1.94
	case km < 50:
		return 1.66
	case km < 100:
		return 1.53
	case km < 5000:
		return 1.41
	default:
		return 1.31
	}
}

func weightFactor(kg float64) float64 {
	switch {
	case kg < 50:
		return 1.88
	case kg < 100:
		return 1.74
	case kg < 250:
		return 1.55
	case kg < 500:
		return 1.37
	default:
		return 1.26
	}
}
