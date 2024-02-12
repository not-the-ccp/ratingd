package tariff

func weightFactor(kg float64) float64 {
	switch {
	case kg < 100:
		return 1.74
	case kg < 500:
		return 1.47
	case kg < 2500:
		return 1.36
	case kg < 5000:
		return 1.20
	case kg < 10000:
		return 1.06
	default:
		return 0.78
	}
}

func zoneFactor(zone float64) float64 {
	switch {
	case zone < 50:
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
