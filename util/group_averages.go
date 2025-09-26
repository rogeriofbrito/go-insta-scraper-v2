package util

// GroupAverages takes a sorted list of integers and a threshold value.
// It groups consecutive numbers where the difference between neighbors
// is less than or equal to the threshold, and returns a slice containing
// the average of each group.
//
// Example:
//
//	nums := []int{1, 2, 3, 20, 21, 22, 40}
//	threshold := 3
//	result := GroupAverages(nums, threshold)
//	// result = [2, 21, 40]
func GroupAverages(nums []int, threshold int) []float64 {
	if len(nums) == 0 {
		return nil
	}

	var result []float64
	sum := nums[0]
	count := 1

	for i := 1; i < len(nums); i++ {
		if nums[i]-nums[i-1] <= threshold {
			// Keep adding to the current group
			sum += nums[i]
			count++
		} else {
			// Close the current group and save its average
			result = append(result, float64(sum)/float64(count))
			// Start a new group
			sum = nums[i]
			count = 1
		}
	}

	// Append the last group's average
	result = append(result, float64(sum)/float64(count))

	return result
}
