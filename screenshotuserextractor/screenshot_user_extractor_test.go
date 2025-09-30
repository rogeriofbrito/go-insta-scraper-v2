package screenshotuserextractor_test

import (
	"image"
	"testing"

	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
	"github.com/rogeriofbrito/go-insta-scraper-v2/screenshotuserextractor"
	"github.com/rogeriofbrito/go-insta-scraper-v2/templatematcher"
	"github.com/rogeriofbrito/go-insta-scraper-v2/tesseractocr"
	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
	"gocv.io/x/gocv"
)

func TestScreenshotUserExtractor_GetUsernames_DiverseCases(t *testing.T) {
	tests := []struct {
		name                  string
		screenshotPath        string
		templateFollowPath    string
		templateFollowingPath string
		templateMessagePath   string
		config                config.Config
		expectedUsernames     []string
		expectErr             bool
	}{
		{
			name:                  "iphone_14_plus_1",
			screenshotPath:        "testdata/iphone_14_plus_1/screenshot.png",
			templateFollowPath:    "testdata/iphone_14_plus_1/follow.png",
			templateFollowingPath: "testdata/iphone_14_plus_1/following.png",
			templateMessagePath:   "testdata/iphone_14_plus_1/following.png", // TODO: change ScreenshotUserExtractor to accept omit templates
			config: config.Config{
				WorkingDirPath:             "/tmp/go-insta-scraper",
				ReferencePointsSearchRect:  image.Rect(600, 308, 675, 1690),
				ReferencePointsXCoordinate: 629,
				GroupAveragesThreshold:     10,
				MatchTemplateThreshold:     float32(0.8),
				MatchTemplateMethod:        gocv.TmCcoeffNormed,
				MatchTemplateImageFlags:    gocv.IMReadColor,
				OcrImageFlags:              gocv.IMReadGrayScale,
				UniformThresold:            5,
				SamplePosition: config.SamplePosition{
					ReferencePoint:        image.Pt(629, 501),
					TopCenterUsernameRect: image.Rect(165, 482, 165+440, 482+36),
					CenterUsernameRect:    image.Rect(165, 518, 165+440, 518+36),
					UpUsernameRect:        image.Rect(165, 498, 165+440, 498+36),
				},
				TesseractOcrOem: 1,
				TesseractOcrPsm: 7, //single text line
				TesseractOcrConfigs: map[string]string{
					"tessedit_char_whitelist":   "abcdefghijklmnopqrstuvwxyz0123456789._",
					"classify_bln_numeric_mode": "1",
					"load_system_dawg":          "0", // disable dictionary corrections
					"load_freq_dawg":            "0", // disable dictionary corrections
				},
			},
			expectedUsernames: []string{
				"matheusgonze1",
				"stephencurry30",
				"siganacaorubronegra",
				"capixabaputo",
				"kvraco",
				"memoriarubronegra",
				"naosalvo",
				"belightstore_",
				"fishfireideas",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := util.CreateWorkingDir(tc.config.WorkingDirPath)
			if err != nil {
				t.Fatalf("error on creating working dir: %v", err)
			}

			tm := templatematcher.NewTemplateMatcher(&tc.config)
			tocr := tesseractocr.NewTesseractOcr(&tc.config)

			extractor := screenshotuserextractor.NewScreenshotUserExtractor(
				tc.screenshotPath,
				tc.templateFollowPath,
				tc.templateFollowingPath,
				tc.templateMessagePath,
				&tc.config,
				tm,
				tocr,
			)

			usernames, err := extractor.GetUsernames()
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !stringSliceEqual(usernames, tc.expectedUsernames) {
				t.Fatalf("GetUsernames() = %v; expected %v", usernames, tc.expectedUsernames)
			}
		})
	}
}

// --- helpers ---

func stringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
