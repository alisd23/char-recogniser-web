package server

import (
	"image"
	"image/png"
	"os"

	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
)

const (
	IMAGE_SIZE = 32
)

// Helper function to save an image to a given path
func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = png.Encode(file, img)
	if err != nil {
		return err
	}
	return nil
}

func NormaliseImage(img image.Image) image.Image {
	resizedImg := resize.Resize(IMAGE_SIZE, IMAGE_SIZE, img, resize.NearestNeighbor)
	grayImg := grayscale.Convert(resizedImg, grayscale.ToGrayLuminance)
	return grayImg
}
