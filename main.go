package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/coltiq/chirpy/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

const (
	port = "8080"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func NewServer(db *database.DB) *http.Server {
	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
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
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpsDelete)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	// Users
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

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
