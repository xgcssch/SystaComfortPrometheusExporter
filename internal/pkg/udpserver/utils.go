package udpserve

func transformToIndicator(indicator int32) float64 {
	if indicator == 0 {
		return float64(0)
	}

	return float64(1)
}
func transformBoolToIndicator(indicator bool) float64 {
	if indicator {
		return float64(0)
	}

	return float64(1)
}
