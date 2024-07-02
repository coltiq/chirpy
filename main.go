package main

import (
	"flag"
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
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	// Users
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return srv
}

func main() {
	databaseName := "database.json"

	db, err := database.NewDB(databaseName)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		log.Print("Debug mode enabled! Deleting database...")
		err := db.ResetDB()
		if err != nil {
			log.Fatalf("Error deleting database: %s", err)
		}
	} else {
		log.Print("Running in normal mode.")
	}

	srv := NewServer(db)

	log.Printf("Starting server [:%s]...", port)
	log.Fatal(srv.ListenAndServe())
}
