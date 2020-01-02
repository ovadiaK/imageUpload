package img

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

const IMAGE_FOLDER_TEMP = "temp-images"
const IMAGE_FOLDER_PERM = "perm-images"
const maxSize = 1400

func Resize(fileName string) (*image.NRGBA, error) {
	var (
		height = maxSize
		width  = maxSize
	)
	i, err := imaging.Open(filepath.Join(IMAGE_FOLDER_TEMP, fileName))
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, fmt.Errorf("i=nil")
	}
	r := i.Bounds()

	if r.Dx() < r.Dy() {
		width = 0
	} else {
		height = 0
	}

	i2 := imaging.Resize(i, width, height, imaging.Lanczos)
	return i2, nil
}

func Format(fileName string) (string, error) {
	src, err := os.Open(filepath.Join(IMAGE_FOLDER_TEMP, fileName))
	if err != nil {
		return fileName, err
	}
	defer src.Close()
	i, _, err := image.Decode(src)
	if err != nil {
		return "", err
	}
	fileName = strings.Join([]string{strings.TrimSuffix(fileName, filepath.Ext(fileName)), "jpg"}, ".")
	dst, err := os.OpenFile(filepath.Join(IMAGE_FOLDER_TEMP, fileName), os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()
	err = jpeg.Encode(dst, i, nil)
	return fileName, err
}
