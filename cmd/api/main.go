package main

import (
    "log"
    "net/http"
    "book-api/internal/routes"

    _ "book-api/docs"
    httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
    router := routes.SetupRoutes()
    router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
    log.Println("Server is RUNNING on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}