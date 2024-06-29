package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(`<html>
				<body>
    				<h1>Welcome, Chirpy Admin</h1>
    				<p>Chirpy has been visited %d times!</p>
				</body>
			</html>`, cfg.fileserverHits)))
}
