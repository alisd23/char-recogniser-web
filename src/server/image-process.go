package server

import (
	"image"
	"image/draw"

	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
)

// NormaliseImage normalises the given image by doing 2 things:
// 1) Resizing to IMAGE_SIZE by IMAGE_SIZE
// 2) Converting to grayscale
func NormaliseImage(img image.Image) image.Image {
	resizedImg := resize.Resize(ImageSize, ImageSize, img, resize.NearestNeighbor)
	grayImg := grayscale.Convert(resizedImg, grayscale.ToGrayLuminance)
	return grayImg
}

// GetPixels extracts the grayscale pixel (intensity) values from an image
func GetPixels(img image.Image) []float32 {
	rect := img.Bounds()
	rawImg := image.NewGray(rect)
	draw.Draw(rawImg, rect, img, rect.Min, draw.Src)
	pixels := rawImg.Pix

	example := make([]float32, len(pixels))

	// Convert into array of intensities - between 0 & 1
	for i, pixel := range pixels {
		example[i] = float32(pixel) / 255
	}
	return example
}
