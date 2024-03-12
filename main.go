package main

import (
	"booking-service/api"
	"booking-service/db"
	"encoding/json"
	_ "fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// HealthCheckResponse represents the structure of the health check JSON response
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	// Connect to the database
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Health check handler function
	healthCheckHandler := func(w http.ResponseWriter, r *http.Request) {
		// Check the database connection
		if err := dbCheck(); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Service is not healthy. Database connection error.")
			return
		}

		// Additional checks can be added here as needed

		// If all checks pass, respond with a success message
		respondWithJSON(w, http.StatusOK, HealthCheckResponse{
			Status:  "UP",
			Message: "Service is up and running!",
		})
	}

	// Register the health check handler function to the "/health" endpoint
	http.HandleFunc("/health", healthCheckHandler)

	// Create a new router
	r := mux.NewRouter()

	// Set up API routes
	api.SetupRoutes(r)

	http.Handle("/", r)

	// Start the web server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// dbCheck checks the database connection status
func dbCheck() error {
	// You can add additional database health checks here if needed
	// For now, just return an error if the connection was not successful
	_, err := db.ConnectDB()
	return err
}

// respondWithError sends a JSON response with an error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
