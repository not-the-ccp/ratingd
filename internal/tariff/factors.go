package tariff

func weightFactor(kg float64) float64 {
	switch {
	case kg < 50:
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
