package server

import (
	"char-recogniser-web/src/database"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/plimble/ace"
	"github.com/rs/cors"
)

const hostname = "localhost"

// Start starts the server
func Start(assetsPath string, port int, pyPort int) {
	router := ace.Default()
	db, err := database.Connect("localhost:27017")

	// Struct containing all endpoint handlers
	endpoints := Endpoints{
		db:         db,
		assetsPath: assetsPath,
		pythonPort: pyPort,
	}

	// HACK Set globally sp this is available in other files

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
	assetsPathAbs, _ := filepath.Abs(filepath.Join(assetsPath, "static"))
	router.Static("/static", http.Dir(assetsPathAbs))

	// API endpoints
	router.GET("/test", func(c *ace.C) { fmt.Println("HERE") })
	router.GET("/api/filters", endpoints.filters)
	router.POST("/api/predict", endpoints.predict)

	// Send index.html on any unmatched url - front-end handles 404
	router.RouteNotFound(endpoints.index)

	url := hostname + ":" + strconv.FormatInt(int64(port), 10)

	handler := c.Handler(router)
	fmt.Printf("SERVER RUNNING ON %#v\n", url)
	err = http.ListenAndServe(url, handler)

	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
