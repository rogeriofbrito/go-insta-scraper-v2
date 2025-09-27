package templatematcher

import (
	"image"
	"image/color"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"gocv.io/x/gocv"
)

// NewTemplateMatcher creates a new TemplateMatcher with the given threshold and matching method.
func NewTemplateMatcher(config *config.Config) *TemplateMatcher {
	return &TemplateMatcher{
		config: config,
	}
}

// TemplateMatcher encapsulates parameters and logic for template matching in images.
type TemplateMatcher struct {
	config *config.Config
}

// GetMatches finds all regions in the image Mat that match the template Mat.
// Returns a slice of rectangles representing the matched regions.
func (tm *TemplateMatcher) GetMatches(imageMat, templateMat gocv.Mat) ([]image.Rectangle, error) {
	// Prepare a result matrix to store match results
	result := gocv.NewMat()
	defer result.Close()

	// Perform template matching
	err := gocv.MatchTemplate(imageMat, templateMat, &result, tm.config.MatchTemplateMethod, gocv.NewMat())
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to match template")
	}

	matches := []image.Rectangle{}
	for {
		// Find the location and value of the best match in the result matrix
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

		// If the best match is below the threshold, stop searching
		if maxVal < tm.config.MatchTemplateThreshold {
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
