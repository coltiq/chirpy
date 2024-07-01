package main

import (
	"log"
	"net/http"

	"github.com/coltiq/chirpy/internal/database"
)

const (
	port = "8080"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func NewServer(db *database.DB) *http.Server {
	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()

	// Wrap Fileserver with Middleware to Track Metrics
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./public")))))

	// Display Server Health
	mux.HandleFunc("GET /api/healthz", HealthHandler)

	// Display Metrics and Reset Metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.ResetMetricsHandler)

	// Create and Retrieve Chirps
	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpsPostHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.ChirpsGetHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.ChirpsGetSingleHandler)

	// Users
	mux.HandleFunc("POST /api/users", apiCfg.UsersPostHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return srv
}

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	srv := NewServer(db)

	log.Printf("Starting server [:%s]...", port)
	log.Fatal(srv.ListenAndServe())
}
