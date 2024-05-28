package main

import (
	"log"
	"net/http"
)

func initServer() {
    const port = "8080"

    mux := http.NewServeMux()

    srv := &http.Server{
                Addr: ":" + port,
                Handler: mux,
    }

    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}

func main() {
    initServer()
}
