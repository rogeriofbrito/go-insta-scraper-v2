package util_test

import (
	"math"
	"testing"

	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
)

// Table-driven tests covering multiple corner cases and typical situations.
func TestGroupAverages_DiverseCases(t *testing.T) {
	tests := []struct {
		name      string
		nums      []int
		threshold int
		expected  []float64 // use nil for expected nil result
	}{
		{
			name:      "empty_input_returns_nil",
			nums:      []int{},
			threshold: 1,
			expected:  nil,
		},
		{
			name:      "single_element",
			nums:      []int{5},
			threshold: 10,
			expected:  []float64{5},
		},
		{
			name:      "all_equal_threshold_zero",
			nums:      []int{2, 2, 2, 2},
			threshold: 0,
			expected:  []float64{2},
		},
		{
			name:      "threshold_zero_distinct_elements_each_own_group",
			nums:      []int{1, 2, 3},
			threshold: 0,
			expected:  []float64{1, 2, 3},
		},
		{
			name:      "large_threshold_groups_everything",
			nums:      []int{1, 2, 3, 4, 5},
			threshold: 100,
			expected:  []float64{3}, // average of 1..5 is 3
		},
		{
			name:      "typical_multiple_groups",
			nums:      []int{1, 2, 3, 20, 21, 22, 40},
			threshold: 3,
			expected:  []float64{2, 21, 40},
		},
		{
			name:      "non_integer_averages",
			nums:      []int{1, 2, 10, 11},
			threshold: 1,
			expected:  []float64{1.5, 10.5},
		},
		{
			name:      "negative_numbers_grouping",
			nums:      []int{-5, -4, -3, 0},
			threshold: 1,
			expected:  []float64{-4, 0},
		},
		{
			name:      "negative_threshold_prevents_grouping",
			nums:      []int{1, 2, 3},
			threshold: -1,
			expected:  []float64{1, 2, 3}, // negative threshold => no adjacent diffs <= threshold
		},
		{
			name:      "groups_with_single_and_multi_element_groups",
			nums:      []int{1, 10, 11, 12, 100},
			threshold: 1,
			expected:  []float64{1, 11, 100}, // groups: [1], [10,11,12], [100]
		},
		{
			name:      "large_numbers_no_overflow_like_case",
			nums:      []int{1 << 30, (1 << 30) + 2},
			threshold: 2,
			// average = ((1<<30) + ((1<<30)+2)) / 2 = 1073741825
			expected: []float64{1073741825.0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.GroupAverages(tc.nums, tc.threshold)
			if !floatSlicesAlmostEqual(got, tc.expected) {
				t.Fatalf("GroupAverages(%v, %d) = %v; expected %v", tc.nums, tc.threshold, got, tc.expected)
			}
		})
	}
}

// helper: compare floats with small tolerance
func floatsAlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

// helper: compare two float64 slices (handles nil)
func floatSlicesAlmostEqual(a, b []float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !floatsAlmostEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
