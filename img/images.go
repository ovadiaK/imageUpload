package img

import (
	"fmt"
	"github.com/disintegration/imaging"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const IMAGE_FOLDER_TEMP = "temp-images"
const IMAGE_FOLDER_PERM = "perm-images"

func Resize(fileName string, size int, makeRectangle bool) (*image.NRGBA, error) {
	width, height := size, size
	i, err := imaging.Open(filepath.Join(IMAGE_FOLDER_TEMP, fileName))
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, fmt.Errorf("i=nil")
	}
	r := i.Bounds()
	if makeRectangle {
		if (r.Dx() / r.Dy()) > 1 {
			dist := (r.Dx() - r.Dy()*2) / 2
			rec := image.Rect(dist, 0, r.Dx()-dist, r.Dy())
			i = imaging.Crop(i, rec)
		}
		if (r.Dy() / r.Dx()) > 1 {
			dist := (r.Dy() - r.Dx()*2) / 2
			rec := image.Rect(0, dist, r.Dx(), r.Dy()-dist)
			i = imaging.Crop(i, rec)
		}
	}
	if r.Dx() < r.Dy() {
		width = 0
	} else {
		height = 0
	}
	i2 := imaging.Resize(i, width, height, imaging.Lanczos)
	return i2, nil
}

func Format(fileName string, format string) (string, error) {
	src, err := os.Open(filepath.Join(IMAGE_FOLDER_TEMP, fileName))
	if err != nil {
		return fileName, err
	}
	defer src.Close()
	i, _, err := image.Decode(src)
	if err != nil {
		return "", err
	}
	fileName = strings.Join([]string{strings.TrimSuffix(fileName, filepath.Ext(fileName)), format}, ".")
	dst, err := os.OpenFile(filepath.Join(IMAGE_FOLDER_TEMP, fileName), os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()
	switch format {
	case "jpg":
		{
			err = jpeg.Encode(dst, i, nil)
		}
	case "png":
		{
			err = png.Encode(dst, i)
		}
	case "tiff":
		{
			err = tiff.Encode(dst, i, nil)
		}
	case "gif":
		{
			err = gif.Encode(dst, i, nil)
		}
	}
	return fileName, err
}
