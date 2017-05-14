package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/plimble/ace"
)

// Default endpoint - serves HTML file
// 404's are handled client-side so any 'unmatched' route uses this handler
func index(c *ace.C) {
	htmlPath, _ := filepath.Abs(filepath.Join(ASSETS_PATH, "index.html"))
	html, err := ioutil.ReadFile(htmlPath)

	if err != nil {
		c.String(500, "Oops, something went wrong")
		fmt.Println(err)
		return
	}

	c.String(200, string(html))
}

type PredictRequest struct {
	Image string
	Test  string
}
type PredictResponse struct {
	Success     bool        `json:"success"`
	Error       string      `json:"error"`
	ShrunkImage image.Image `json:"shrunkImage"`
}

// Predict endpoint
// User sends the letter iamge to predict (in bytes? form)
func predict(c *ace.C) {
	// Unmarshal body
	body := PredictRequest{}
	c.ParseJSON(&body)

	// Convert base 64 image into []byte
	b64data := body.Image[strings.IndexByte(body.Image, ',')+1:]
	imgData, err := base64.StdEncoding.DecodeString(b64data)

	// Decode []byte into image.Image type
	decodedImg, err := png.Decode(bytes.NewReader(imgData))

	if err != nil {
		fmt.Println("ERROR: ", err)
		c.JSON(200, PredictResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	normalisedImg := NormaliseImage(decodedImg)

	err = SaveImage(normalisedImg, "test_image.png")

	c.JSON(200, PredictResponse{
		Success: true,
	})
}
