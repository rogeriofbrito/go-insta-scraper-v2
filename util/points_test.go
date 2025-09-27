package util_test

import (
	"image"
	"testing"

	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
)

func TestGetMinPointsFromRects_DiverseCases(t *testing.T) {
	tests := []struct {
		name     string
		rects    []image.Rectangle
		expected []image.Point
	}{
		{
			name:     "empty_input_returns_empty",
			rects:    nil,
			expected: nil,
		},
		{
			name:     "single_rect",
			rects:    []image.Rectangle{{Min: image.Pt(1, 2), Max: image.Pt(3, 4)}},
			expected: []image.Point{{1, 2}},
		},
		{
			name: "multiple_rects",
			rects: []image.Rectangle{
				{Min: image.Pt(0, 0), Max: image.Pt(1, 1)},
				{Min: image.Pt(-1, -2), Max: image.Pt(5, 5)},
				{Min: image.Pt(10, 20), Max: image.Pt(30, 40)},
			},
			expected: []image.Point{{0, 0}, {-1, -2}, {10, 20}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.GetMinPointsFromRects(tc.rects)
			if !pointsSliceEqual(got, tc.expected) {
				t.Fatalf("GetMinPointsFromRects(%v) = %v; expected %v", tc.rects, got, tc.expected)
			}
		})
	}
}

func TestGetPointsInsideRect_DiverseCases(t *testing.T) {
	tests := []struct {
		name     string
		points   []image.Point
		rect     image.Rectangle
		expected []image.Point
	}{
		{
			name:     "empty_points_returns_empty",
			points:   nil,
			rect:     image.Rect(0, 0, 10, 10),
			expected: nil,
		},
		{
			name:     "all_points_inside",
			points:   []image.Point{{1, 1}, {5, 5}, {9, 9}},
			rect:     image.Rect(0, 0, 10, 10),
			expected: []image.Point{{1, 1}, {5, 5}, {9, 9}},
		},
		{
			name:     "all_points_outside",
			points:   []image.Point{{-1, -1}, {11, 11}, {0, 11}},
			rect:     image.Rect(0, 0, 10, 10),
			expected: nil,
		},
		{
			name:     "some_points_inside_some_outside",
			points:   []image.Point{{0, 0}, {10, 10}, {5, 5}, {11, 0}},
			rect:     image.Rect(0, 0, 10, 10),
			expected: []image.Point{{0, 0}, {10, 10}, {5, 5}},
		},
		{
			name:     "points_on_edges",
			points:   []image.Point{{0, 0}, {10, 10}, {0, 10}, {10, 0}},
			rect:     image.Rect(0, 0, 10, 10),
			expected: []image.Point{{0, 0}, {10, 10}, {0, 10}, {10, 0}},
		},
		{
			name:     "rect_with_negative_coords",
			points:   []image.Point{{-5, -5}, {-1, -1}, {0, 0}},
			rect:     image.Rect(-5, -5, 0, 0),
			expected: []image.Point{{-5, -5}, {-1, -1}, {0, 0}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.GetPointsInsideRect(tc.points, tc.rect)
			if !pointsSliceEqual(got, tc.expected) {
				t.Fatalf("GetPointsInsideRect(%v, %v) = %v; expected %v", tc.points, tc.rect, got, tc.expected)
			}
		})
	}
}

func TestGetYCoordinatesFromPoints_DiverseCases(t *testing.T) {
	tests := []struct {
		name     string
		points   []image.Point
		expected []int
	}{
		{
			name:     "empty_points_returns_empty",
			points:   nil,
			expected: nil,
		},
		{
			name:     "single_point",
			points:   []image.Point{{3, 7}},
			expected: []int{7},
		},
		{
			name:     "multiple_points",
			points:   []image.Point{{1, 2}, {3, 4}, {5, 6}},
			expected: []int{2, 4, 6},
		},
		{
			name:     "negative_y_values",
			points:   []image.Point{{0, -1}, {0, -10}},
			expected: []int{-1, -10},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.GetYCoordinatesFromPoints(tc.points)
			if !intSliceEqual(got, tc.expected) {
				t.Fatalf("GetYCoordinatesFromPoints(%v) = %v; expected %v", tc.points, got, tc.expected)
			}
		})
	}
}

func TestGetReferencePoints_DiverseCases(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		ys       []int
		expected []image.Point
	}{
		{
			name:     "empty_ys_returns_empty",
			x:        5,
			ys:       nil,
			expected: nil,
		},
		{
			name:     "single_y",
			x:        2,
			ys:       []int{7},
			expected: []image.Point{{2, 7}},
		},
		{
			name:     "multiple_ys",
			x:        0,
			ys:       []int{1, 2, 3},
			expected: []image.Point{{0, 1}, {0, 2}, {0, 3}},
		},
		{
			name:     "negative_x_and_ys",
			x:        -3,
			ys:       []int{-1, -2},
			expected: []image.Point{{-3, -1}, {-3, -2}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.GetReferencePoints(tc.x, tc.ys)
			if !pointsSliceEqual(got, tc.expected) {
				t.Fatalf("GetReferencePoints(%d, %v) = %v; expected %v", tc.x, tc.ys, got, tc.expected)
			}
		})
	}
}

// --- helpers ---

func pointsSliceEqual(a, b []image.Point) bool {
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

func intSliceEqual(a, b []int) bool {
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
