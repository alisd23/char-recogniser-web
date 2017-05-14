package server

import (
	"char-recogniser-go/src/database"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	mgo "gopkg.in/mgo.v2"

	"github.com/plimble/ace"
	"github.com/rs/cors"
)

const HOSTNAME = "localhost"
const PORT = 9000

var ASSETS_PATH string
var DB *mgo.Database

func Start(assetsPath string) {
	ASSETS_PATH = assetsPath

	router := ace.Default()
	connection, err := database.Connect("localhost:27017")

	// HACK Set globally sp this is available in other files
	DB = connection

	if err != nil {
		fmt.Println("[SERVER] Database connection error: ", err)
		return
	}

	// Allow requests from ANY localhost
	// NOTE: Should turn this off in production
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowCredentials: true,
	})

	// Server static files
	assetsPathAbs, _ := filepath.Abs(filepath.Join(ASSETS_PATH, "static"))
	router.Static("/static", http.Dir(assetsPathAbs))

	// API endpoints
	router.POST("/api/predict", predict)
	router.POST("/api/train", train)

	// Send index.html on any unmatched url - front-end handles 404
	router.RouteNotFound(index)

	url := HOSTNAME + ":" + strconv.FormatInt(PORT, 10)
	fmt.Printf("SERVER RUNNING ON %#v\n", url)

	handler := c.Handler(router)
	err = http.ListenAndServe(url, handler)

	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
