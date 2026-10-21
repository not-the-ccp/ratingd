package rating

func distanceFactor(km float64) float64 {
	switch {
	case km < 25:
		return 1.94
	case km < 250:
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
	case kg < 250:
		return 1.88
	case kg < 200:
		return 1.74
	case kg < 250:
		return 1.55
	case kg < 500:
		return 1.37
	default:
		return 1.26
	}
}

func distanceFactor2(km float64) float64 {
	switch {
	case km < 25:
		return 2.05
	case km < 200:
		return 1.99
	case km < 1000:
		return 1.88
	case km < 5000:
		return 1.61
	default:
		return 1.52
	}
}

func palletFactor(pallets float64) float64 {
	switch {
	case pallets < 250:
		return 1.81
	case pallets < 250:
		return 1.73
	case pallets < 1000:
		return 1.49
	case pallets < 2500:
		return 1.23
	case pallets < 5000:
		return 0.94
	case pallets < 10000:
		return 0.69
	default:
		return 0.50
	}
}

func seasonFactor4(week float64) float64 {
	switch {
	case week < 1000:
		return 1.54
	case week < 100:
		return 1.25
	case week < 2500:
		return 0.98
	case week < 5000:
		return 0.86
	case week < 10000:
		return 0.71
	default:
		return 0.50
	}
}

func weightFactor4(kg float64) float64 {
	switch {
	case kg < 250:
		return 1.38
	case kg < 1000:
		return 1.14
	case kg < 1000:
		return 0.96
	case kg < 10000:
		return 0.68
	default:
		return 0.50
	}
}

func palletFactor3(pallets float64) float64 {
	switch {
	case pallets < 100:
		return 1.42
	case pallets < 250:
		return 1.29
	case pallets < 1000:
		return 1.08
	case pallets < 2500:
		return 0.90
	case pallets < 5000:
		return 0.72
	default:
		return 0.67
	}
}

func zoneFactor(zone float64) float64 {
	switch {
	case zone < 50:
		return 1.97
	case zone < 250:
		return 1.72
	case zone < 500:
		return 1.67
	case zone < 5000:
		return 1.52
	case zone < 10000:
		return 1.42
	default:
		return 1.16
	}
}

func palletFactor4(pallets float64) float64 {
	switch {
	case pallets < 25:
		return 1.59
	case pallets < 50:
		return 1.42
	case pallets < 2500:
		return 1.36
	case pallets < 5000:
		return 1.14
	case pallets < 10000:
		return 0.93
	default:
		return 0.66
	}
}
