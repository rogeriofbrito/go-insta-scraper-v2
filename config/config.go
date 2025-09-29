package config

import (
	"image"

	"gocv.io/x/gocv"
)

type SamplePosition struct {
	ReferencePoint        image.Point     // Min point of a rectangle that surrounds a button
	CenterUsernameRect    image.Rectangle // Center username rectangle relative to reference point
	TopCenterUsernameRect image.Rectangle // Top Center username rectangle relative to reference point
	UpUsernameRect        image.Rectangle // Up username rectangle relative to reference point
}

type Config struct {
	WorkingDirPath                         string                 // Path tp directory when temp files will be created
	ReferencePointsSearchRect              image.Rectangle        // Area where reference points are allowed to be in (changes according device/fontsize where image was captured)
	ReferencePointsXCoordinate             int                    // X coordinate of reference points
	GroupAveragesThreshold                 int                    // Maximum difference between two consecutive numbers for them to belong to the same group
	MatchTemplateThreshold                 float32                // Minimum similarity threshold for a match to be considered valid
	MatchTemplateMethod                    gocv.TemplateMatchMode // Template matching method (e.g., SQDIFF, CCORR)
	ScreenshotUserExtractorImageFlags      gocv.IMReadFlag        // Flags used in Screenshot User Extractor to read screenshot and template images
	ScreenshotUserExtractorUniformThresold int                    // Maximum difference threshold between pixel values for them to be considered equal
	SamplePosition                         SamplePosition         // Sample of a reference point ant 3 rectangles, that will be used to define base rectangles
	TesseractOcrOem                        int                    // Tesseract OCR engine mode (OEM) to use for text recognition
	TesseractOcrPsm                        int                    // Tesseract OCR page segmentation mode (PSM) to use for text recognition
	TesseractOcrConfigs                    map[string]string      // Additional Tesseract OCR configuration key-value pairs
}
