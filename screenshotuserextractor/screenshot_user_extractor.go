package screenshotuserextractor

import (
	"image"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
	"gocv.io/x/gocv"
)

func NewScreenshotUserExtractor(
	screenshotMat gocv.Mat,
	templateFollowMat gocv.Mat,
	templateFollowingMat gocv.Mat,
	templateMessageMat gocv.Mat,
	config *config.Config,
	tm *templatematcher.TemplateMatcher,
) *ScreenshotUserExtractor {
	return &ScreenshotUserExtractor{
		screenshotMat:        screenshotMat,
		templateFollowMat:    templateFollowMat,
		templateFollowingMat: templateFollowingMat,
		templateMessageMat:   templateFollowMat,
		config:               config,
		tm:                   tm,
	}
}

type ScreenshotUserExtractor struct {
	screenshotMat        gocv.Mat
	templateFollowMat    gocv.Mat
	templateFollowingMat gocv.Mat
	templateMessageMat   gocv.Mat
	config               *config.Config
	tm                   *templatematcher.TemplateMatcher
}

func (s *ScreenshotUserExtractor) GetUsernames() ([]string, error) {
	matches, err := s.getMatches()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get matches")
	}

	minPoints := util.GetMinPointsFromRects(matches)
	minPointsSecure := util.GetPointsInsideRect(minPoints, s.config.ReferencePointsSearchRect)
	yCoordinates := util.GetYCoordinatesFromPoints(minPointsSecure)
	yCoordinatesGroup := util.GroupAverages(yCoordinates, s.config.GroupAveragesThreshold)
	yCoordinatesGroupInt := util.ConvertSliceFloat64ToInt(yCoordinatesGroup)
	referencePoints := util.GetReferencePoints(s.config.ReferencePointsXCoordinate, yCoordinatesGroupInt)

	var usernameRects []image.Rectangle
	for _, referencePoint := range referencePoints {
		topCenterUsernameRect := s.getTopCenterUsernameRect(referencePoint)
		if util.IsUniformRegion(s.screenshotMat, topCenterUsernameRect, 5) { //TODO: create config for threshold
			usernameRects = append(usernameRects, s.getCenterUsernameRect(referencePoint))
		} else {
			usernameRects = append(usernameRects, s.getUpUsernameRect(referencePoint))
		}
	}

	// TODO: extract usernames with tesseract

	return nil, nil
}

func (s *ScreenshotUserExtractor) getMatches() ([]image.Rectangle, error) {
	var matches []image.Rectangle

	matchesFollow, err := s.tm.GetMatches(s.screenshotMat, s.templateFollowMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get follow buttom matches")
	}

	matchesFollowing, err := s.tm.GetMatches(s.screenshotMat, s.templateFollowingMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get following buttom matches")
	}

	matchesMessage, err := s.tm.GetMatches(s.screenshotMat, s.templateMessageMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get message buttom matches")
	}

	matches = append(matches, matchesFollow...)
	matches = append(matches, matchesFollowing...)
	matches = append(matches, matchesMessage...)

	return matches, nil
}

func (s *ScreenshotUserExtractor) getCenterUsernameRect(referencePoint image.Point) image.Rectangle {
	return image.Rect(referencePoint.X-465, referencePoint.Y+15, referencePoint.X-135, referencePoint.Y+56) // TODO: move values to config
}

func (s *ScreenshotUserExtractor) getTopCenterUsernameRect(referencePoint image.Point) image.Rectangle {
	return image.Rect(referencePoint.X-465, referencePoint.Y+5, referencePoint.X-135, referencePoint.Y+18) // TODO: move values to config
}

func (s *ScreenshotUserExtractor) getUpUsernameRect(referencePoint image.Point) image.Rectangle {
	return image.Rect(referencePoint.X-465, referencePoint.Y-3, referencePoint.X-135, referencePoint.Y+34) // TODO: move values to config
}
