package main

import (
	"log"
	"net/http"
)

const (
	port = "8080"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) ResetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func NewServer() *http.Server {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()

	// Wrap Fileserver with Middleware to Track Metrics
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./public")))))

	// Display Server Health
	mux.HandleFunc("GET /api/healthz", HealthHandler)

	// Display Metrics and Reset Metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.ResetMetricsHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return srv
}

func main() {
	srv := NewServer()

	log.Printf("Starting server [:%s]...", port)
	log.Fatal(srv.ListenAndServe())
}
