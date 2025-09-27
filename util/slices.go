package util

func ConvertSliceFloat64ToInt(floatSlice []float64) []int {
	var intSlice []int
	for _, floatValue := range floatSlice {
		intSlice = append(intSlice, int(floatValue))
	}

	return intSlice
}
