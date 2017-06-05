package server

import (
	"bytes"
	"char-recogniser-go/src/database"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"path/filepath"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"github.com/plimble/ace"
)

func isValidCharCode(charCode int) bool {
	return ((charCode >= 48 && charCode <= 57) ||
		(charCode >= 65 && charCode <= 90) ||
		(charCode >= 97 && charCode <= 122))
}

// Struct containing all endpoint handlers
type Endpoints struct {
	db         *mgo.Database
	assetsPath string
}

// Default endpoint - serves HTML file
// 404's are handled client-side so any 'unmatched' route uses this handler
func (e Endpoints) index(c *ace.C) {
	htmlPath, _ := filepath.Abs(filepath.Join(e.assetsPath, "index.html"))
	html, err := ioutil.ReadFile(htmlPath)

	if err != nil {
		c.String(500, "Oops, something went wrong")
		fmt.Println(err)
		return
	}

	c.String(200, string(html))
}

type TrainRequest struct {
	Image    string
	CharCode int
}
type TrainResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Train endpoint
// User sends image and expected char code
func (e Endpoints) train(c *ace.C) {
	// Unmarshal body
	fmt.Println("Handling train request")
	body := TrainRequest{}
	c.ParseJSON(&body)

	// Convert base 64 image into []byte, trimming off type information first
	b64data := body.Image[strings.IndexByte(body.Image, ',')+1:]
	imgData, err := base64.StdEncoding.DecodeString(b64data)

	// Decode []byte into image.Image type
	decodedImg, err := png.Decode(bytes.NewReader(imgData))

	if err != nil || !isValidCharCode(body.CharCode) {
		fmt.Println("[SERVER] /api/train ERROR: ", err)
		c.JSON(200, TrainResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	normalisedImg := NormaliseImage(decodedImg)

	record := database.CreateExample(normalisedImg, body.CharCode)
	database.InsertExample(e.db, record)

	c.JSON(200, PredictResponse{
		Success: true,
	})
}

type PredictRequest struct {
	Image string
}
type PredictResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Predict endpoint
// User sends the letter iamge to predict (in bytes? form)
func (e Endpoints) predict(c *ace.C) {
	// Unmarshal body
	body := PredictRequest{}
	c.ParseJSON(&body)

	// Convert base 64 image into []byte, trimming off type information first
	b64data := body.Image[strings.IndexByte(body.Image, ',')+1:]
	imgData, err := base64.StdEncoding.DecodeString(b64data)

	// Decode []byte into image.Image type
	decodedImg, err := png.Decode(bytes.NewReader(imgData))

	if err != nil {
		fmt.Println("[SERVER] /api/predict ERROR: ", err)
		c.JSON(200, PredictResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	normalisedImg := NormaliseImage(decodedImg)

	fmt.Println(normalisedImg.Bounds())

	c.JSON(200, PredictResponse{
		Success: true,
	})
}
