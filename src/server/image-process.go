package server

import (
	"image"

	"char-recogniser-go/src/database"

	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
)

func NormaliseImage(img image.Image) image.Image {
	resizedImg := resize.Resize(database.IMAGE_SIZE, database.IMAGE_SIZE, img, resize.NearestNeighbor)
	grayImg := grayscale.Convert(resizedImg, grayscale.ToGrayLuminance)
	return grayImg
}
