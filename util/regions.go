package util

import (
	"image"
	"math"

	"gocv.io/x/gocv"
)

// IsUniformRegion checks if all pixels in a given rectangular region of an image
// are similar within a specified threshold.
// imageMat: the input image as a gocv.Mat
// rect: the region of interest as an image.Rectangle
// threshold: the maximum allowed difference between pixel values
// Returns true if the region is uniform, false otherwise.
func IsUniformRegion(imageMat gocv.Mat, rect image.Rectangle, threshold int) bool {
	// Extract the region of interest from the image.
	region := imageMat.Region(rect)
	defer region.Close() // Ensure the region is released after use.

	// Get the value of the first pixel in the region as a reference.
	refPixel := region.GetVecbAt(0, 0)

	// If the image is single-channel (grayscale)
	if imageMat.Channels() == 1 {
		firstPixel := refPixel[0]

		// Iterate over all pixels in the region.
		for y := 0; y < region.Rows(); y++ {
			for x := 0; x < region.Cols(); x++ {
				currentPixel := region.GetUCharAt(y, x)
				// If the difference exceeds the threshold, region is not uniform.
				if math.Abs(float64(currentPixel)-float64(firstPixel)) > float64(threshold) {
					return false
				}
			}
		}
		return true // All pixels are within the threshold.
	} else {
		// For multi-channel images (e.g., BGR color)
		firstPixel := refPixel
		b := firstPixel[0]
		g := firstPixel[1]
		r := firstPixel[2]

		// Iterate over all pixels in the region.
		for y := 0; y < region.Rows(); y++ {
			for x := 0; x < region.Cols(); x++ {
				currentPixel := region.GetVecbAt(y, x)
				cb := currentPixel[0]
				cg := currentPixel[1]
				cr := currentPixel[2]
				// Check if any channel exceeds the threshold.
				if math.Abs(float64(cb)-float64(b)) > float64(threshold) ||
					math.Abs(float64(cg)-float64(g)) > float64(threshold) ||
					math.Abs(float64(cr)-float64(r)) > float64(threshold) {
					return false
				}
			}
		}
		return true // All pixels are within the threshold for all channels.
	}
}
