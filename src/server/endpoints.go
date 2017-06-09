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
	res, _ := http.Post(pyPredictURL, "application/json; charset=utf-8", b)

	var pythonRes pythonPredictRes

	json.NewDecoder(res.Body).Decode(&pythonRes)

	if pythonRes.Predictions == nil {
		fmt.Println("Python response invalid", pythonRes)
		c.JSON(500, predictResponse{
			Error: "Prediction request failed",
		})
		return
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, normalisedImg)
	imageBytes := buf.Bytes()

	if err != nil {
		fmt.Println("Image encoding failed", err)
		c.JSON(500, predictResponse{
			Error: "Prediction request failed",
		})
		return
	}
	c.JSON(200, predictResponse{
		Predictions: pythonRes.Predictions,
		Image:       base64.StdEncoding.EncodeToString(imageBytes),
	})
}

// Parameters endpoint
// Fetches various model parameters (e.g. conv filters)
func (e Endpoints) parameters(c *ace.C) {
	var result map[string]interface{}

	fs := e.db.GridFS("training_runs")
	err := fs.Find(nil).Sort("_id").One(&result)

	if err != nil {
		fmt.Println("[/parameters] ERROR finding GridFS file ", err)
		return
	}

	file, err := fs.OpenId(result["_id"])

	if err != nil {
		fmt.Println("[/parameters] ERROR opening GridFS file", err)
		return
	}

	var buf bytes.Buffer
	buf.ReadFrom(file)
	buf.WriteTo(c.Writer)
}
