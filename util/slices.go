package util

func ConvertSliceFloat64ToInt(floatSlice []float64) []int {
	var intSlice []int
	for _, floatValue := range floatSlice {
		intSlice = append(intSlice, int(floatValue))
	}

	return intSlice
}

func RemoveEmptyString(strs []string) []string {
	if strs == nil {
		return nil
	}

	strsNonEmpty := []string{}
	for _, str := range strs {
		if len(str) > 0 {
			strsNonEmpty = append(strsNonEmpty, str)
		}
	}

	return strsNonEmpty
}
