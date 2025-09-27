package util_test

import (
	"image"
	"testing"

	"gocv.io/x/gocv"

	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
)

// Table-driven tests covering multiple corner cases and typical situations.
func TestIsUniformRegion_DiverseCases(t *testing.T) {
	type scenario struct {
		name      string
		imagePath string
		flags     gocv.IMReadFlag
		threshold int
		expected  bool
	}

	cases := []scenario{
		{
			name:      "non_uniform_image_1",
			imagePath: "testdata/points/non_uniform_images/image_1.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  false,
		},
		{
			name:      "non_uniform_image_2",
			imagePath: "testdata/points/non_uniform_images/image_2.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  false,
		},
		{
			name:      "non_uniform_image_3",
			imagePath: "testdata/points/non_uniform_images/image_3.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  false,
		},
		{
			name:      "non_uniform_image_4",
			imagePath: "testdata/points/non_uniform_images/image_4.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  false,
		},
		{
			name:      "non_uniform_image_5",
			imagePath: "testdata/points/non_uniform_images/image_5.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  false,
		},
		{
			name:      "uniform_image_1",
			imagePath: "testdata/points/uniform_images/image_1.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  true,
		},
		{
			name:      "uniform_image_2",
			imagePath: "testdata/points/uniform_images/image_2.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  true,
		},
		{
			name:      "uniform_image_3",
			imagePath: "testdata/points/uniform_images/image_3.png",
			flags:     gocv.IMReadColor,
			threshold: 5,
			expected:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			imageMat := gocv.IMRead(tc.imagePath, tc.flags)
			if imageMat.Empty() {
				t.Fatalf("failed to load image: %s", tc.imagePath)
			}
			defer imageMat.Close()

			rect := image.Rect(0, 0, imageMat.Cols()-1, imageMat.Rows()-1)
			got := util.IsUniformRegion(imageMat, rect, tc.threshold)
			if got != tc.expected {
				t.Errorf("IsUniformRegion(%s, %v, %d) = %v; want %v", tc.imagePath, rect, tc.threshold, got, tc.expected)
			}
		})
	}
}
