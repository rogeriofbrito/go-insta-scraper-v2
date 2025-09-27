package main

import (
	"image"

	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"github.com/rogeriofbrito/go-insta-scraper-v2/screenshotuserextractor"
	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"gocv.io/x/gocv"
)

func main() {
	//TODO: read screenshot and templates inside ScreenshotUserExtractor

	screenshotMat := gocv.IMRead("./frame/frame_0056.png", gocv.IMReadColor)
	if screenshotMat.Empty() {
		panic("failed to read screenshot")
	}
	defer screenshotMat.Close()

	templateFollowMat := gocv.IMRead("./template/pt_BR/follow.png", gocv.IMReadColor)
	if templateFollowMat.Empty() {
		panic("failed to read follow buttom template")
	}
	defer templateFollowMat.Close()

	templateFollowingMat := gocv.IMRead("./template/pt_BR/following.png", gocv.IMReadColor)
	if templateFollowingMat.Empty() {
		panic("failed to read following buttom template")
	}
	defer templateFollowingMat.Close()

	templateMessageMat := gocv.IMRead("./template/pt_BR/message.png", gocv.IMReadColor)
	if templateMessageMat.Empty() {
		panic("failed to read message buttom template")
	}
	defer templateMessageMat.Close()

	// TODO: add function to create config.Config from env vars
	config := &config.Config{
		ReferencePointsSearchRect:  image.Rect(630, 280, 650, 1800),
		ReferencePointsXCoordinate: 635,
		GroupAveragesThreshold:     10,
		MatchTemplateThreshold:     float32(0.8),
		MatchTemplateMethod:        gocv.TmCcoeffNormed,
	}
	tm := templatematcher.NewTemplateMatcher(config)
	sue := screenshotuserextractor.NewScreenshotUserExtractor(screenshotMat, templateFollowMat, templateFollowingMat, templateMessageMat, config, tm)

	_, err := sue.GetUsernames()
	if err != nil {
		panic(err)
	}
}
