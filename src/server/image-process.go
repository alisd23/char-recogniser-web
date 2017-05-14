package server

import (
	"image"
	"image/png"
	"os"

	"char-recogniser-go/src/database"

	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
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
	resizedImg := resize.Resize(database.IMAGE_SIZE, database.IMAGE_SIZE, img, resize.NearestNeighbor)
	grayImg := grayscale.Convert(resizedImg, grayscale.ToGrayLuminance)
	return grayImg
}
