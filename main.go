package main

import (
	"fmt"
	"net/http"
)

func createServer() {
    handler := http.NewServeMux()
    server := http.Server{
                Addr: "localhost:8080",
                Handler: handler,
    }
    err := server.ListenAndServe()
    if err != nil {
        fmt.Println("bad")
    }
}

func main() {
    createServer()
}
