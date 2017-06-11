package server

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"

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

// PixelsToBase64 converts an array of pixels to its base64 encoding
func PixelsToBase64(pixels []float32, stride int) (string, error) {
	intPixels := make([]uint8, stride*stride)
	for j, p := range pixels {
		intPixels[j] = uint8(p)
	}

	img := &image.Gray{
		Pix:    intPixels,
		Stride: stride,
		Rect:   image.Rect(0, 0, stride, stride),
	}
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)

	if err != nil {
		return "", err
	}

	imageBytes := buf.Bytes()

	base64 := base64.StdEncoding.EncodeToString(imageBytes)

	return base64, nil
}
