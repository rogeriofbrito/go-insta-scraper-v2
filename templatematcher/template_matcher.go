package templatematcher

import (
	"image"
	"image/color"

	"github.com/palantir/stacktrace"
	"gocv.io/x/gocv"
)

func NewTemplateMatcher(threshold float32, flags gocv.IMReadFlag, method gocv.TemplateMatchMode) *TemplateMatcher {
	return &TemplateMatcher{
		threshold: threshold,
		flags:     flags,
		method:    method,
	}
}

type TemplateMatcher struct {
	threshold float32
	flags     gocv.IMReadFlag
	method    gocv.TemplateMatchMode
}

func (tm *TemplateMatcher) GetMatches(imagePath, templatePath string) ([]*image.Rectangle, error) {
	imageMat := gocv.IMRead(imagePath, tm.flags)
	if imageMat.Empty() {
		return nil, stacktrace.NewError("failed to read image at %s", imagePath)
	}
	defer imageMat.Close()

	templateMat := gocv.IMRead(templatePath, tm.flags)
	if templateMat.Empty() {
		return nil, stacktrace.NewError("failed to read image at %s", templatePath)
	}
	defer templateMat.Close()

	result := gocv.NewMat()
	defer result.Close()

	err := gocv.MatchTemplate(imageMat, templateMat, &result, tm.method, gocv.NewMat())
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to match template")
	}

	matches := []*image.Rectangle{}
	for {
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

		if maxVal < tm.threshold {
			break
		}

		x0, y0, x1, y1 := maxLoc.X, maxLoc.Y, maxLoc.X+templateMat.Cols(), maxLoc.Y+templateMat.Rows()
		match := image.Rect(x0, y0, x1, y1)

		err := gocv.Rectangle(&result, match, color.RGBA{0, 0, 0, 0}, -1)
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed to suppress match")
		}

		matches = append(matches, &match)
	}

	return matches, nil
}
