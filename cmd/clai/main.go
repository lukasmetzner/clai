package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/lukasmetzner/clai/internal/clai/handlers"
	"github.com/lukasmetzner/clai/pkg/database"
	"github.com/lukasmetzner/clai/pkg/mq"
)

// FileServer custom handler to serve the SPA
func spaHandler(staticPath, indexPath string) http.Handler {
	mime.AddExtensionType(".js", "application/javascript")
	fs := http.FileServer(http.Dir(staticPath))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the absolute path to the requested file
		path, err := filepath.Abs(filepath.Join(staticPath, r.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the file exists and is not a directory
		stat, err := os.Stat(path)
		if os.IsNotExist(err) || stat.IsDir() {
			// If the file does not exist or is a directory, serve the index.html
			http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serve the file
		fs.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize the database
	database.InitDB()

	// Initialize RabbitMQ Client
	mq.InitMQ()

	// Set up the router
	r := mux.NewRouter()

	// Define the endpoints
	r.HandleFunc("/api/jobs", handlers.CreateJob).Methods("POST")
	r.HandleFunc("/api/jobs", handlers.GetJobs).Methods("GET")
	r.HandleFunc("/api/jobs/{id}", handlers.GetJob).Methods("GET")
	r.HandleFunc("/api/jobs/{id}", handlers.UpdateJob).Methods("PUT")
	r.HandleFunc("/api/jobs/{id}", handlers.DeleteJob).Methods("DELETE")

	r.PathPrefix("/").Handler(spaHandler("./frontend/dist", "index.html"))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
