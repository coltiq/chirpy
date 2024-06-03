package main

import (
    "fmt"
	"log"
	"net/http"
)

type apiConfig struct {
    fileserverHits  int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        cfg.fileserverHits++
        next.ServeHTTP(w, r)
    })
}

func initServer() {
    const filepathRoot = "."
    const port = ":8080"
    apiCfg := &apiConfig{
            fileserverHits: 0,
    }

    mux := http.NewServeMux()

    // Register a custom Handler
    http.HandleFunc("/healthz", healthHandler)

    fileServer := http.FileServer(http.Dir(filepathRoot))

    // Metrics Handler
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
    mux.HandleFunc("/metrics", apiCfg.hitsHandler)
    mux.HandleFunc("/reset", apiCfg.resetHandler)

    log.Fatal(http.ListenAndServe(port, mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits)))
}

func main() {
    initServer()
}
