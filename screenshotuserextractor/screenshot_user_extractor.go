package screenshotuserextractor

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"github.com/rogeriofbrito/go-insta-scraper-v2/tesseractocr"
	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
	"gocv.io/x/gocv"
)

func NewScreenshotUserExtractor(
	screenshotPath string,
	templateFollowPath string,
	templateFollowingPath string,
	templateMessagePath string,
	config *config.Config,
	tm *templatematcher.TemplateMatcher,
	tocr *tesseractocr.TesseractOcr,
) *ScreenshotUserExtractor {
	return &ScreenshotUserExtractor{
		screenshotPath:        screenshotPath,
		templateFollowPath:    templateFollowPath,
		templateFollowingPath: templateFollowingPath,
		templateMessagePath:   templateMessagePath,
		config:                config,
		tm:                    tm,
		tocr:                  tocr,
	}
}

type ScreenshotUserExtractor struct {
	screenshotPath        string
	templateFollowPath    string
	templateFollowingPath string
	templateMessagePath   string
	config                *config.Config
	tm                    *templatematcher.TemplateMatcher
	tocr                  *tesseractocr.TesseractOcr
}

func (s *ScreenshotUserExtractor) GetUsernames() ([]string, error) {
	screenshotMat, err := s.readImage(s.screenshotPath, s.config.ScreenshotUserExtractorImageFlags)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read screenshot image")
	}

	templateFollowMat, err := s.readImage(s.templateFollowPath, s.config.ScreenshotUserExtractorImageFlags)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read follow template image")
	}

	templateFollowingMat, err := s.readImage(s.templateFollowingPath, s.config.ScreenshotUserExtractorImageFlags)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read following template image")
	}

	templateMessageMat, err := s.readImage(s.templateMessagePath, s.config.ScreenshotUserExtractorImageFlags)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read message template image")
	}

	defer screenshotMat.Close()
	defer templateFollowMat.Close()
	defer templateFollowingMat.Close()
	defer templateMessageMat.Close()

	matches, err := s.getMatches(screenshotMat, templateFollowMat, templateFollowingMat, templateMessageMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get matches")
	}

	minPoints := util.GetMinPointsFromRects(matches)
	minPointsSecure := util.GetPointsInsideRect(minPoints, s.config.ReferencePointsSearchRect)
	yCoordinates := util.GetYCoordinatesFromPoints(minPointsSecure)
	yCoordinatesGroup := util.GroupAverages(yCoordinates, s.config.GroupAveragesThreshold)
	yCoordinatesGroupInt := util.ConvertSliceFloat64ToInt(yCoordinatesGroup)
	referencePoints := util.GetReferencePoints(s.config.ReferencePointsXCoordinate, yCoordinatesGroupInt)
	usernameRects := s.getUsernameRects(screenshotMat, referencePoints)

	usernameImagePaths, err := s.writeUsernameImages(screenshotMat, usernameRects)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to write username images")
	}

	usernames, err := s.ocrUsernames(usernameImagePaths)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read usernames from screenshot")
	}

	return usernames, nil
}

func (s *ScreenshotUserExtractor) readImage(imagePath string, flags gocv.IMReadFlag) (gocv.Mat, error) {
	imageMat := gocv.IMRead(imagePath, flags)
	if imageMat.Empty() {
		return gocv.Mat{}, stacktrace.NewError("failed to read image at %s: image empty", imagePath)
	}

	return imageMat, nil
}

func (s *ScreenshotUserExtractor) getMatches(
	screenshotMat,
	templateFollowMat,
	templateFollowingMat,
	templateMessageMat gocv.Mat,
) ([]image.Rectangle, error) {
	var matches []image.Rectangle

	matchesFollow, err := s.tm.GetMatches(screenshotMat, templateFollowMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get follow buttom matches")
	}

	matchesFollowing, err := s.tm.GetMatches(screenshotMat, templateFollowingMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get following buttom matches")
	}

	matchesMessage, err := s.tm.GetMatches(screenshotMat, templateMessageMat)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get message buttom matches")
	}

	matches = append(matches, matchesFollow...)
	matches = append(matches, matchesFollowing...)
	matches = append(matches, matchesMessage...)

	return matches, nil
}

func (s *ScreenshotUserExtractor) getUsernameRects(screenshotMat gocv.Mat, referencePoints []image.Point) []image.Rectangle {
	baseTopCenterUsernameRect := s.config.SamplePosition.TopCenterUsernameRect.Sub(s.config.SamplePosition.ReferencePoint)
	baseCenterUsernameRect := s.config.SamplePosition.CenterUsernameRect.Sub(s.config.SamplePosition.ReferencePoint)
	baseUpUsernameRect := s.config.SamplePosition.UpUsernameRect.Sub(s.config.SamplePosition.ReferencePoint)

	var usernameRects []image.Rectangle
	for _, referencePoint := range referencePoints {
		topCenterUsernameRect := baseTopCenterUsernameRect.Add(referencePoint)
		if util.IsUniformRegion(screenshotMat, topCenterUsernameRect, s.config.ScreenshotUserExtractorUniformThresold) {
			usernameRects = append(usernameRects, baseCenterUsernameRect.Add(referencePoint))
		} else {
			usernameRects = append(usernameRects, baseUpUsernameRect.Add(referencePoint))
		}
	}

	return usernameRects
}

func (s *ScreenshotUserExtractor) writeUsernameImages(screenshotMat gocv.Mat, usernameRects []image.Rectangle) ([]string, error) {
	var usernameImagePaths []string
	for i, usernameRect := range usernameRects {
		usernameMat := screenshotMat.Region(usernameRect)
		usernameImagePath := fmt.Sprintf("%s/username_%d.jpg", s.config.WorkingDirPath, i)

		writeSuccess := gocv.IMWrite(usernameImagePath, usernameMat)
		if !writeSuccess {
			return nil, stacktrace.NewError("failed to write mat at path %s", usernameImagePath)
		}

		usernameImagePaths = append(usernameImagePaths, usernameImagePath)
	}

	return usernameImagePaths, nil
}

func (s *ScreenshotUserExtractor) ocrUsernames(usernameImagePaths []string) ([]string, error) {
	var usernames []string
	for _, usernameImagePath := range usernameImagePaths {
		err := s.tocr.OCR(usernameImagePath, usernameImagePath)
		if err != nil {
			return nil, stacktrace.Propagate(err,
				"failed to execute tesseract ocr over %s with result at %s",
				usernameImagePath, usernameImagePath+".txt")
		}

		usernameOcrTxtBytes, err := os.ReadFile(usernameImagePath + ".txt")
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed to read file %s", usernameImagePath+".txt")
		}

		usernameOcrTxt := string(usernameOcrTxtBytes)
		usernameOcrTxtLines := strings.Split(usernameOcrTxt, "\n")
		usernameOcrTxtLines = util.RemoveEmptyString(usernameOcrTxtLines)

		if len(usernameOcrTxtLines) != 1 {
			return nil, stacktrace.NewError(
				"failed to execute ocr in image %s: number of lines returned is different than 1 (%d)",
				usernameImagePath,
				len(usernameOcrTxtLines),
			)
		}

		usernames = append(usernames, usernameOcrTxtLines[0])
	}

	return usernames, nil
}
