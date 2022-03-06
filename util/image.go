package util

import (
	"bytes"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

func ValidateImageType(fileType string) bool {
	switch strings.ToLower(fileType) {
	case "image/png":
		return true
	case "image/jpg":
		return true
	case "image/jpeg":
		return true
	case "image/webp":
		return true
	//Todo 處理未上傳檔案的情況
	case "":
		return true
	default:
		return false
	}
}

func DecodeBase64AndSaveAsWebp(base64EncodedData string, dst string) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if err := webp.Encode(out, img, options); err != nil {
		return err
	}

	return nil
}