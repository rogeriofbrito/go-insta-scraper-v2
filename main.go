package main

import (
	"image"

	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"github.com/rogeriofbrito/go-insta-scraper-v2/screenshotuserextractor"
	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"gocv.io/x/gocv"
)

func main() {
	// TODO: add function to create config.Config from env vars
	config := &config.Config{
		ReferencePointsSearchRect:              image.Rect(630, 280, 650, 1800),
		ReferencePointsXCoordinate:             635,
		GroupAveragesThreshold:                 10,
		MatchTemplateThreshold:                 float32(0.8),
		MatchTemplateMethod:                    gocv.TmCcoeffNormed,
		ScreenshotUserExtractorImageFlags:      gocv.IMReadColor,
		ScreenshotUserExtractorUniformThresold: 5,
		BaseCenterUsernameRect:                 image.Rect(-465, 15, -135, 56),
		BaseTopCenterUsernameRect:              image.Rect(-465, 5, -135, 18),
		BaseUpUsernameRect:                     image.Rect(-465, -3, -135, 34),
	}
	tm := templatematcher.NewTemplateMatcher(config)
	sue := screenshotuserextractor.NewScreenshotUserExtractor(
		"./frame/frame_0056.png",
		"./template/pt_BR/follow.png",
		"./template/pt_BR/following.png",
		"./template/pt_BR/message.png",
		config,
		tm,
	)

	_, err := sue.GetUsernames()
	if err != nil {
		panic(err)
	}
}
