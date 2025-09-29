package main

import (
	"image"

	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"github.com/rogeriofbrito/go-insta-scraper-v2/screenshotuserextractor"
	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"github.com/rogeriofbrito/go-insta-scraper-v2/tesseractocr"
	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
	"gocv.io/x/gocv"
)

func main() {
	// TODO: add function to create config.Config from env vars
	config := &config.Config{
		WorkingDirPath:                         "/tmp/go-insta-scraper",
		ReferencePointsSearchRect:              image.Rect(600, 308, 675, 1690),
		ReferencePointsXCoordinate:             629,
		GroupAveragesThreshold:                 10,
		MatchTemplateThreshold:                 float32(0.8),
		MatchTemplateMethod:                    gocv.TmCcoeffNormed,
		ScreenshotUserExtractorImageFlags:      gocv.IMReadColor,
		ScreenshotUserExtractorUniformThresold: 5,
		SamplePosition: config.SamplePosition{
			ReferencePoint:        image.Pt(629, 501),
			TopCenterUsernameRect: image.Rect(165, 482, 165+440, 482+36),
			CenterUsernameRect:    image.Rect(165, 518, 165+440, 518+36),
			UpUsernameRect:        image.Rect(165, 498, 165+440, 498+36),
		},
		TesseractOcrOem: 1,
		TesseractOcrPsm: 7,
		TesseractOcrConfigs: map[string]string{
			"tessedit_char_whitelist":   "abcdefghijklmnopqrstuvwxyz0123456789._",
			"classify_bln_numeric_mode": "1",
		},
	}

	err := util.CreateWorkingDir(config.WorkingDirPath)
	if err != nil {
		panic(err)
	}

	tm := templatematcher.NewTemplateMatcher(config)
	tocr := tesseractocr.NewTesseractOcr(config)

	sue := screenshotuserextractor.NewScreenshotUserExtractor(
		"./frame/frame_0056.png",         //TODO: move to config
		"./template/pt_BR/follow.png",    //TODO: move to config
		"./template/pt_BR/following.png", //TODO: move to config
		"./template/pt_BR/message.png",   //TODO: move to config
		config,
		tm,
		tocr,
	)

	_, err = sue.GetUsernames()
	if err != nil {
		panic(err)
	}
}
