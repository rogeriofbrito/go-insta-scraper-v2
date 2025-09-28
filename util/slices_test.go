package util_test

import (
	"testing"

	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
)

// Table-driven tests covering multiple corner cases and typical situations.
func TestConvertSliceFloat64ToInt_DiverseCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []int
	}{
		{
			name:     "empty_input_returns_empty",
			input:    nil,
			expected: nil,
		},
		{
			name:     "single_zero",
			input:    []float64{0},
			expected: []int{0},
		},
		{
			name:     "single_positive_integer",
			input:    []float64{42},
			expected: []int{42},
		},
		{
			name:     "single_negative_integer",
			input:    []float64{-17},
			expected: []int{-17},
		},
		{
			name:     "multiple_integers",
			input:    []float64{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "floats_truncated_towards_zero",
			input:    []float64{1.9, 2.1, -3.7, -4.2, 0.5},
			expected: []int{1, 2, -3, -4, 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.ConvertSliceFloat64ToInt(tc.input)
			if !intSlicesEqual(got, tc.expected) {
				t.Fatalf("ConvertSliceFloat64ToInt(%v) = %v; expected %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestRemoveEmptyString_DiverseCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty_input_returns_empty",
			input:    nil,
			expected: nil,
		},
		{
			name:     "all_empty_strings",
			input:    []string{"", "", ""},
			expected: []string{},
		},
		{
			name:     "no_empty_strings",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "mixed_empty_and_nonempty",
			input:    []string{"", "foo", "", "bar", "", "baz"},
			expected: []string{"foo", "bar", "baz"},
		},
		{
			name:     "single_nonempty",
			input:    []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "single_empty",
			input:    []string{""},
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.RemoveEmptyString(tc.input)
			if !stringSlicesEqual(got, tc.expected) {
				t.Fatalf("RemoveEmptyString(%v) = %v; expected %v", tc.input, got, tc.expected)
			}
		})
	}
}

// helper: compare two int slices (handles nil)
func intSlicesEqual(a, b []int) bool {
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
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// helper: compare two string slices (handles nil)
func stringSlicesEqual(a, b []string) bool {
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
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
