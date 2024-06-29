package main

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
