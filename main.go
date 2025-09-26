package main

import (
	"fmt"

	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"gocv.io/x/gocv"
)

func main() {
	tm := templatematcher.NewTemplateMatcher(float32(0.8), gocv.IMReadColor, gocv.TmCcoeffNormed)
	matches, err := tm.GetMatches("./frame/frame_0056.png", "./template/pt_BR/follow.png")
	if err != nil {
		panic(err)
	}

	fmt.Println(matches)
}
