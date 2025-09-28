package tesseractocr

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/rogeriofbrito/go-insta-scraper-v2/config"
)

func NewTesseractOcr(config *config.Config) *TesseractOcr {
	return &TesseractOcr{
		oem:     config.TesseractOcrOem,
		psm:     config.TesseractOcrPsm,
		configs: config.TesseractOcrConfigs,
	}
}

type TesseractOcr struct {
	oem     int
	psm     int
	configs map[string]string
}

func (t *TesseractOcr) OCR(imagePath string, resultPath string) error {
	args := []string{
		imagePath,
		resultPath,
		"--oem",
		strconv.Itoa(t.oem),
		"--psm",
		strconv.Itoa(t.psm),
	}
	args = append(args, t.getConfigArgs()...)

	cmd := exec.Command("tesseract", args...)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (t *TesseractOcr) getConfigArgs() []string {
	var configArgs []string
	for configName, configValue := range t.configs {
		configArgs = append(configArgs, "-c")
		configArgs = append(configArgs, fmt.Sprintf("%s=%s", configName, configValue))
	}

	return configArgs
}
