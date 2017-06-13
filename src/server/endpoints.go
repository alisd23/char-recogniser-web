package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"github.com/plimble/ace"
)

func isValidCharCode(charCode int) bool {
	return ((charCode >= 48 && charCode <= 57) ||
		(charCode >= 65 && charCode <= 90) ||
		(charCode >= 97 && charCode <= 122))
}

// Endpoints is a struct containing all endpoint handlers
type Endpoints struct {
	db         *mgo.Database
	assetsPath string
	pythonPort int
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

type predictRequest struct {
	Image string
}
type predictResponse struct {
	Error       string       `json:"error"`
	Image       string       `json:"image"`
	Predictions []prediction `json:"predictions"`
	Activations []string     `json:"activations"`
}
type prediction struct {
	Charcode   int    `json:"charcode"`
	Confidence string `json:"confidence"`
}
type pythonPredictReq struct {
	Image []float32 `json:"image"`
}
type pythonPredictRes struct {
	Predictions []prediction `json:"predictions"`
	Activations [][]float32  `json:"activations"`
}

// Predict endpoint
// User sends the letter iamge to predict (in bytes? form)
func (e Endpoints) predict(c *ace.C) {
	// Unmarshal body
	body := predictRequest{}
	c.ParseJSON(&body)

	// Convert base 64 image into []byte, trimming off type information first
	b64data := body.Image[strings.IndexByte(body.Image, ',')+1:]
	imgData, err := base64.StdEncoding.DecodeString(b64data)

	// Decode []byte into image.Image type
	decodedImg, err := png.Decode(bytes.NewReader(imgData))

	if err != nil {
		fmt.Println("[SERVER] /api/predict ERROR: ", err)
		c.JSON(500, predictResponse{
			Error: err.Error(),
		})
		return
	}

	normalisedImg := NormaliseImage(decodedImg)

	pixels := GetPixels(normalisedImg)

	pythonBody := pythonPredictReq{
		Image: pixels,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(pythonBody)
	pyPredictURL := "http://localhost:" + strconv.Itoa(e.pythonPort) + "/predict"
	res, err := http.Post(pyPredictURL, "application/json; charset=utf-8", b)

	if err != nil {
		fmt.Println("Python response failed")
		c.JSON(500, predictResponse{
			Error: "Prediction request failed",
		})
		return
	}

	var pythonRes pythonPredictRes

	json.NewDecoder(res.Body).Decode(&pythonRes)

	if pythonRes.Predictions == nil {
		fmt.Println("Python response invalid", pythonRes)
		c.JSON(500, predictResponse{
			Error: "Prediction request failed",
		})
		return
	}

	for i, p := range pixels {
		pixels[i] = p * 255
	}
	imgString, err := PixelsToBase64(pixels, ImageSize)

	if err != nil {
		fmt.Println("Image encoding failed", err)
		c.JSON(500, predictResponse{
			Error: "Prediction request failed",
		})
	}

	activationImgs := make([]string, noOfFilters)

	for i, filter := range pythonRes.Activations {
		imgBase64, err := PixelsToBase64(filter, ImageSize)
		if err == nil {
			activationImgs[i] = imgBase64
		}
	}

	c.JSON(200, predictResponse{
		Predictions: pythonRes.Predictions,
		Activations: activationImgs,
		Image:       imgString,
	})
}

type modelValues struct {
	Top1  float32
	Top3  float32
	Conv1 [][]float32
}
type filtersResponse struct {
	Filters []string `json:"filters"`
	Top1    float32  `json:"top1"`
	Top3    float32  `json:"top3"`
}

const noOfFilters = 32
const filterSize = 5

// Model endpoint
// Fetches the filters in image format for the conv layer 1
func (e Endpoints) model(c *ace.C) {
	var result modelValues
	err := e.db.C("values").Find(nil).One(&result)

	if err != nil {
		fmt.Println("[ERROR] No values found")
		c.JSON(500, predictResponse{
			Error: "No values found",
		})
		return
	}

	images := make([]string, noOfFilters)

	for i, filter := range result.Conv1 {
		img, err := PixelsToBase64(filter, filterSize)
		if err == nil {
			images[i] = img
		}
	}

	c.JSON(200, filtersResponse{
		Filters: images,
		Top1:    result.Top1,
		Top3:    result.Top3,
	})
}
