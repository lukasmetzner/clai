package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/lukasmetzner/clai/internal/clai/handlers"
	"github.com/lukasmetzner/clai/pkg/database"
	"github.com/lukasmetzner/clai/pkg/mq"
)

func spaHandler(staticPath, indexPath string) gin.HandlerFunc {
	mime.AddExtensionType(".js", "application/javascript")
	fs := http.FileServer(http.Dir(staticPath))

	return func(c *gin.Context) {
		// Get the absolute path to the requested file
		path, err := filepath.Abs(filepath.Join(staticPath, c.Request.URL.Path))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		// Check if the file exists and is not a directory
		stat, err := os.Stat(path)
		if os.IsNotExist(err) || stat.IsDir() {
			// If the file does not exist or is a directory, serve the index.html
			c.File(filepath.Join(staticPath, indexPath))
			return
		} else if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Serve the file
		fs.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// Initialize the database
	database.InitDB()

	// Initialize RabbitMQ Client
	mq.InitMQ()

	// Set up the router
	r := gin.Default()

	api := r.Group("/api")

	jobs := api.Group("/jobs")
	jobs.POST("/", handlers.CreateJob)
	jobs.GET("/", handlers.GetJobs)
	jobs.GET("/:id", handlers.GetJob)
	jobs.PUT("/:id", handlers.UpdateJob)
	jobs.DELETE("/:id", handlers.DeleteJob)

	// Serve Frontend
	r.GET("/", spaHandler("./frontend/dist", "index.html"))
	r.NoRoute(spaHandler("./frontend/dist", "index.html"))

	// Start the server
	log.Println("Starting server on :8080")
	r.Run()
}
