package main

import (
    "log"
    "net/http"
    "book-api/internal/routes"
)

func main() {
    router := routes.SetupRoutes()
	log.Println("Server is RUNNING on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}