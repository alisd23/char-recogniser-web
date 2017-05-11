package server

import (
	"image"

	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
)

const (
	IMAGE_SIZE = 32
)

func normaliseImage(img image.Image) image.Image {
	resizedImg := resize.Resize(IMAGE_SIZE, IMAGE_SIZE, img, resize.NearestNeighbor)
	grayImg := grayscale.Convert(resizedImg, grayscale.ToGrayLuminance)
	return grayImg
}
