package util

import (
	"bytes"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
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

func GetImageType(fileType string) string {
	return strings.Split(fileType, "/")[1]
}

func DecodeBase64AndSaveAsWebp(base64EncodedData string, imageType string, dst string) (int, int, error) {
	//dst = dst + ".webp"
	//options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)
	//if err != nil {
	//	return 0, 0, err
	//}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return 0, 0, err
	}

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
			Quality: 100,
		})
		if err != nil {
			return 0, 0, err
		}

	case "image/jpeg":
		dst = dst + ".jpeg"
		f, _ := os.Create(dst)
		defer f.Close()
		err := png.Encode(f, img)
		if err != nil {
			return 0, 0, err
		}

	default:
		return 0, 0, nil
	}

	//out, err := os.Create(dst)
	//if err != nil {
	//	return 0, 0, err
	//}

	//if err := webp.Encode(out, img, options); err != nil {
	//	return 0, 0, err
	//}

	return width, height, nil
}

func DecodeBase64AndUpdateAsWebp(base64EncodedData string, imageType string, dst string) error {
	//options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)
	//if err != nil {
	//	return err
	//}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(dst, os.O_RDWR|os.O_TRUNC, 0666)
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
			Quality: 100,
		})
		if err != nil {
			return err
		}

	case "image/jpeg":
		err := png.Encode(f, img)
		if err != nil {
			return err
		}

	default:
		return nil
	}

	//if err := webp.Encode(out, img, options); err != nil {
	//	return err
	//}

	return nil
}

func DecodeBase64AndSaveAsWebpInLQIP(base64EncodedData string, imageType string, dst string) error {
	//dst = dst + "_lq.webp"
	//options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	//if err != nil {
	//	return err
	//}

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

	//out, err := os.Create(dst)
	//if err != nil {
	//	return err
	//}
	//
	//if err := webp.Encode(out, blurImage, options); err != nil {
	//	return err
	//}

	return nil
}

func DecodeBase64AndUpdateAsWebpInLQIP(base64EncodedData string, imageType string, dst string) error {
	//options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	//if err != nil {
	//	return err
	//}

	img, _, err := image.Decode(bytes.NewReader([]byte(base64EncodedData)))
	if err != nil {
		return err
	}

	src := imaging.Resize(img, 200, 0, imaging.NearestNeighbor)

	blurImage := imaging.Blur(src, 5)

	f, err := os.OpenFile(dst, os.O_RDWR|os.O_TRUNC, 0666)
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

	//if err := webp.Encode(out, blurImage, options); err != nil {
	//	return err
	//}

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
