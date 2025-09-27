package util

import "image"

// GetMinPointsFromRects returns the minimum (top-left) points of each rectangle in the input slice.
func GetMinPointsFromRects(rects []image.Rectangle) []image.Point {
	var minPoints []image.Point
	for _, rect := range rects {
		minPoints = append(minPoints, rect.Min)
	}

	return minPoints
}

// GetPointsInsideRect filters and returns only the points that are inside the given rectangle.
func GetPointsInsideRect(points []image.Point, rect image.Rectangle) []image.Point {
	var pointsIn []image.Point
	for _, point := range points {
		if pointInRect(point, rect) {
			pointsIn = append(pointsIn, point)
		}
	}

	return pointsIn
}

// GetYCoordinatesFromPoints extracts and returns the Y coordinates from a slice of points.
func GetYCoordinatesFromPoints(points []image.Point) []int {
	var yCoordinates []int
	for _, point := range points {
		yCoordinates = append(yCoordinates, point.Y)
	}

	return yCoordinates
}

// GetReferencePoints generates reference points used to locate other elements in an image.
// All reference points share the same X coordinate (x), and their Y coordinates are provided by the ys slice.
func GetReferencePoints(x int, ys []int) []image.Point {
	var referencePoints []image.Point
	for _, y := range ys {
		refPoint := image.Pt(x, y)
		referencePoints = append(referencePoints, refPoint)
	}

	return referencePoints
}

// pointInRect checks if a given point is inside the specified rectangle.
func pointInRect(point image.Point, rect image.Rectangle) bool {
	return point.X >= rect.Min.X &&
		point.X <= rect.Max.X &&
		point.Y >= rect.Min.Y &&
		point.Y <= rect.Max.Y
}
