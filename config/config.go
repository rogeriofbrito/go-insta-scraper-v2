package config

import (
	"image"

	"gocv.io/x/gocv"
)

type Config struct {
	WorkingDirPath                         string                 // Path tp directory when temp files will be created
	ReferencePointsSearchRect              image.Rectangle        // Area where reference points are allowed to be in (changes according device/fontsize where image was captured)
	ReferencePointsXCoordinate             int                    // X coordinate of reference points
	GroupAveragesThreshold                 int                    // Maximum difference between two consecutive numbers for them to belong to the same group
	MatchTemplateThreshold                 float32                // Minimum similarity threshold for a match to be considered valid
	MatchTemplateMethod                    gocv.TemplateMatchMode // Template matching method (e.g., SQDIFF, CCORR)
	ScreenshotUserExtractorImageFlags      gocv.IMReadFlag        // Flags used in Screenshot User Extractor to read screenshot and template images
	ScreenshotUserExtractorUniformThresold int                    // Maximum difference threshold between pixel values for them to be considered equal
	BaseCenterUsernameRect                 image.Rectangle        // Rectangle that is translated by reference point to generate CenterUsernameRect
	BaseTopCenterUsernameRect              image.Rectangle        // Rectangle that is translated by reference point to generate TopCenterUsernameRect
	BaseUpUsernameRect                     image.Rectangle        // Rectangle that is translated by reference point to generate UpUsernameRect
	TesseractOcrOem                        int                    // Tesseract OCR engine mode (OEM) to use for text recognition
	TesseractOcrPsm                        int                    // Tesseract OCR page segmentation mode (PSM) to use for text recognition
	TesseractOcrConfigs                    map[string]string      // Additional Tesseract OCR configuration key-value pairs
}
