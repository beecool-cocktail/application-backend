package util

import (
	"bytes"
	"fmt"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/disintegration/imaging"
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

func DecodeBase64AndSaveAsWebp(base64EncodedData string, dst string) (int, int, error) {
	dst = dst + ".webp"
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)
	if err != nil {
		return 0, 0, err
	}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return 0, 0, err
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	out, err := os.Create(dst)
	if err != nil {
		return 0, 0, err
	}

	if err := webp.Encode(out, img, options); err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func DecodeBase64AndUpdateAsWebp(base64EncodedData string, dst string) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	if err := webp.Encode(out, img, options); err != nil {
		return err
	}

	return nil
}

func DecodeBase64AndSaveAsWebpInLQIP(base64EncodedData string, dst string) error {
	dst = dst + "_lq.webp"
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	src := imaging.Resize(img, 200, 0, imaging.NearestNeighbor)

	blurImage := imaging.Blur(src, 5)
	fmt.Println(dst)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if err := webp.Encode(out, blurImage, options); err != nil {
		return err
	}

	return nil
}

func DecodeBase64AndUpdateAsWebpInLQIP(base64EncodedData string, dst string) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	src := imaging.Resize(img, 200, 0, imaging.NearestNeighbor)

	blurImage := imaging.Blur(src, 5)

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	if err := webp.Encode(out, blurImage, options); err != nil {
		return err
	}

	return nil
}

func GetFileNameByPath(path string) (string, error) {
	pathSplitArray := strings.Split(path, "/")

	//目前db的path為 static/fileName
	if len(pathSplitArray) == 2 {
		return pathSplitArray[1], nil
	} else {
		return "", domain.ErrFilePathIllegal
	}
}
