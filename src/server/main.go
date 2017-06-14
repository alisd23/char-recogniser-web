package server

import (
	"char-recogniser-web/src/database"
	"fmt"
	"net/http"
	"os"

	"github.com/plimble/ace"
	"github.com/rs/cors"
)

const defaultServerHost = "localhost"
const defaultDBHost = "localhost"
const defaultPythonHost = "localhost"
const defaultPythonPort = "9001"

// Start starts the server
func Start(assetsPath string) {
	router := ace.Default()

	serverHost := os.Getenv("SERVER_HOST")
	dbHost := os.Getenv("DB_HOST")
	pythonHost := os.Getenv("PYTHON_HOST")
	pythonPort := os.Getenv("PYTHON_PORT")

	if serverHost == "" {
		serverHost = defaultServerHost
	}
	if dbHost == "" {
		dbHost = defaultDBHost
	}
	if pythonHost == "" {
		pythonHost = defaultPythonHost
	}
	if pythonPort == "" {
		pythonPort = defaultPythonPort
	}

	db, err := database.Connect(dbHost + ":27017")

	// Struct containing all endpoint handlers
	endpoints := Endpoints{
		db:         db,
		assetsPath: assetsPath,
		pythonHost: pythonHost,
		pythonPort: pythonPort,
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
	router.Static("/build", http.Dir(assetsPath))

	// API endpoints
	router.GET("/api/model", endpoints.model)
	router.POST("/api/predict", endpoints.predict)

	// Send index.html on any unmatched url - front-end handles 404
	router.RouteNotFound(endpoints.index)

	url := serverHost + ":9000"

	handler := c.Handler(router)
	fmt.Printf("SERVER RUNNING ON %#v\n", url)
	err = http.ListenAndServe(url, handler)

	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
