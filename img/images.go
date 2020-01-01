package img

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"mime/multipart"
)

const IMAGE_FOLDER = "temp-images"
const maxSize = 1400

func Resize(file *multipart.File) (*image.NRGBA, error) {
	var (
		height = maxSize
		width  = maxSize
	)
	conf, format, err := image.DecodeConfig(*file)
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(*file)
	if err != nil {
		return nil, err
	}
	fmt.Println(format)

	if conf.Height < conf.Width {
		width = 0
	} else {
		height = 0
	}

	i2 := imaging.Resize(i, width, height, imaging.Lanczos)
	return i2, nil
}
