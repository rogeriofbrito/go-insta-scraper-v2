package config

import (
	"image"

	"gocv.io/x/gocv"
)

type Config struct {
	ReferencePointsSearchRect         image.Rectangle        // Area where reference points are allowed to be in (changes according device/fontsize where image was captured)
	ReferencePointsXCoordinate        int                    // X coordinate of reference points
	GroupAveragesThreshold            int                    // Minimum difference between two consecutive numbers for them to belong to the same group
	MatchTemplateThreshold            float32                // Minimum similarity threshold for a match to be considered valid
	MatchTemplateMethod               gocv.TemplateMatchMode // Template matching method (e.g., SQDIFF, CCORR)
	ScreenshotUserExtractorImageFlags gocv.IMReadFlag        // Flags used in Screenshot User Extractor to read screenshot and template images
	BaseCenterUsernameRect            image.Rectangle        // Rectangle that is translated by reference point to generate CenterUsernameRect
	BaseTopCenterUsernameRect         image.Rectangle        // Rectangle that is translated by reference point to generate TopCenterUsernameRect
	BaseUpUsernameRect                image.Rectangle        // Rectangle that is translated by reference point to generate UpUsernameRect
}
