package util

import (
	"encoding/base64"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
	"mime/multipart"
	"os"
	"strings"
	_ "image/jpeg"
	_ "image/png"
)

func DecodeImage(data string) (image []byte, fileType string, err error) {
	// imageData: data:image/<format>;base64,<data>
	dataSplit := strings.Split(data, ",")

	// dataSplit[0]: data:image/<format>;base64
	fileType = strings.Split(strings.Replace(dataSplit[0], ";", "/", -1), "/")[1]

	image, err = base64.StdEncoding.DecodeString(dataSplit[1])
	if err != nil {
		return
	}

	return
}

func ValidateImageType(fileType string) bool {
	switch strings.ToLower(fileType) {
	case ".png":
		return true
	case ".jpg":
		return true
	case ".webp":
		return true
	default:
		return false
	}
}

func SaveAsWebp(fileHeader *multipart.FileHeader, dst string) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
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