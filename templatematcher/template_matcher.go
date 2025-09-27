package templatematcher

import (
	"image"
	"image/color"

	"github.com/palantir/stacktrace"
	"gocv.io/x/gocv"
)

// NewTemplateMatcher creates a new TemplateMatcher with the given threshold, image read flags, and matching method.
func NewTemplateMatcher(threshold float32, flags gocv.IMReadFlag, method gocv.TemplateMatchMode) *TemplateMatcher {
	return &TemplateMatcher{
		threshold: threshold,
		flags:     flags,
		method:    method,
	}
}

// TemplateMatcher encapsulates parameters and logic for template matching in images.
type TemplateMatcher struct {
	threshold float32                // Minimum similarity threshold for a match to be considered valid
	flags     gocv.IMReadFlag        // Flags for reading images (e.g., color, grayscale)
	method    gocv.TemplateMatchMode // Template matching method (e.g., SQDIFF, CCORR)
}

// GetMatches finds all regions in the image at imagePath that match the template at templatePath.
// Returns a slice of rectangles representing the matched regions.
func (tm *TemplateMatcher) GetMatches(imagePath, templatePath string) ([]image.Rectangle, error) {
	// Read the main image from disk
	imageMat := gocv.IMRead(imagePath, tm.flags)
	if imageMat.Empty() {
		return nil, stacktrace.NewError("failed to read image at %s", imagePath)
	}
	defer imageMat.Close()

	// Read the template image from disk
	templateMat := gocv.IMRead(templatePath, tm.flags)
	if templateMat.Empty() {
		return nil, stacktrace.NewError("failed to read image at %s", templatePath)
	}
	defer templateMat.Close()

	// Prepare a result matrix to store match results
	result := gocv.NewMat()
	defer result.Close()

	// Perform template matching
	err := gocv.MatchTemplate(imageMat, templateMat, &result, tm.method, gocv.NewMat())
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to match template")
	}

	matches := []image.Rectangle{}
	for {
		// Find the location and value of the best match in the result matrix
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

		// If the best match is below the threshold, stop searching
		if maxVal < tm.threshold {
			break
		}

		// Define the rectangle for the matched region
		x0, y0, x1, y1 := maxLoc.X, maxLoc.Y, maxLoc.X+templateMat.Cols(), maxLoc.Y+templateMat.Rows()
		match := image.Rect(x0, y0, x1, y1)

		// Suppress this match in the result matrix to avoid duplicate detections
		err := gocv.Rectangle(&result, match, color.RGBA{0, 0, 0, 0}, -1)
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed to suppress match")
		}

		// Add the matched rectangle to the results
		matches = append(matches, match)
	}

	return matches, nil
}
