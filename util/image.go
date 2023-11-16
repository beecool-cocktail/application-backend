package util

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/disintegration/imaging"
	"github.com/vincent-petithory/dataurl"
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

func GetImageType(fileType string) string {
	return strings.Split(fileType, "/")[1]
}

func DecodeBase64AndSaveAsWebp(base64EncodedData string, imageType string, dst string) (int, int, error) {

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return 0, 0, err
	}

	compressionRatio := getCompressionRatio(base64EncodedData, domain.AllowMaxImageSizeInMB)

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	switch imageType {
	case "image/png":
		dst = dst + ".png"
		f, _ := os.Create(dst)
		defer f.Close()
		err := png.Encode(f, img)
		if err != nil {
			return 0, 0, err
		}

	case "image/jpg":
		dst = dst + ".jpg"
		f, _ := os.Create(dst)
		defer f.Close()
		err := jpeg.Encode(f, img, &jpeg.Options{
			Quality: compressionRatio,
		})
		if err != nil {
			return 0, 0, err
		}

	case "image/jpeg":
		dst = dst + ".jpeg"
		f, _ := os.Create(dst)
		defer f.Close()
		err := jpeg.Encode(f, img, &jpeg.Options{
			Quality: compressionRatio,
		})
		if err != nil {
			return 0, 0, err
		}

	default:
		return 0, 0, nil
	}

	return width, height, nil
}

func DecodeBase64AndUpdateAsWebp(base64EncodedData string, imageType string, sourceFile, destinationFile string) error {

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	compressionRatio := getCompressionRatio(base64EncodedData, domain.AllowMaxImageSizeInMB)

	f, err := os.OpenFile(sourceFile, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	switch imageType {
	case "image/png":
		err := png.Encode(f, img)
		if err != nil {
			return err
		}

	case "image/jpg":
		err := jpeg.Encode(f, img, &jpeg.Options{
			Quality: compressionRatio,
		})
		if err != nil {
			return err
		}

	case "image/jpeg":
		err := jpeg.Encode(f, img, &jpeg.Options{
			Quality: compressionRatio,
		})
		if err != nil {
			return err
		}

	default:
		return nil
	}

	err = os.Rename(sourceFile, destinationFile)
	if err != nil {
		return err
	}

	return nil
}

func DecodeBase64AndSaveAsWebpInLQIP(base64EncodedData string, imageType string, dst string) error {

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	src := imaging.Resize(img, 200, 0, imaging.NearestNeighbor)

	blurImage := imaging.Blur(src, 5)

	switch imageType {
	case "image/png":
		dst = dst + "_lq.png"
		f, _ := os.Create(dst)
		defer f.Close()
		err := png.Encode(f, blurImage)
		if err != nil {
			return err
		}

	case "image/jpg":
		dst = dst + "_lq.jpg"
		f, _ := os.Create(dst)
		defer f.Close()
		err := jpeg.Encode(f, blurImage, &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return err
		}

	case "image/jpeg":
		dst = dst + "_lq.jpeg"
		f, _ := os.Create(dst)
		defer f.Close()
		err := png.Encode(f, blurImage)
		if err != nil {
			return err
		}

	default:
		return nil
	}

	return nil
}

func DecodeBase64AndUpdateAsWebpInLQIP(base64EncodedData string, imageType string, sourceFile, destinationFile string) error {

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	src := imaging.Resize(img, 200, 0, imaging.NearestNeighbor)

	blurImage := imaging.Blur(src, 5)

	f, err := os.OpenFile(sourceFile, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	switch imageType {
	case "image/png":
		err := png.Encode(f, blurImage)
		if err != nil {
			return err
		}

	case "image/jpg":
		err := jpeg.Encode(f, blurImage, &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return err
		}

	case "image/jpeg":
		err := png.Encode(f, blurImage)
		if err != nil {
			return err
		}

	default:
		return nil
	}

	err = os.Rename(sourceFile, destinationFile)
	if err != nil {
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

func ParseLQIPFileToDataURL(file string) (string, error) {

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	dataURL := dataurl.EncodeBytes(bytes)

	return dataURL, nil
}

func getCompressionRatio(image string, targetSizeInMB int) int {
	originSizeInMB := len(image) / (1024 * 1024)
	if originSizeInMB < 5 {
		return 100
	}

	return originSizeInMB / targetSizeInMB
}