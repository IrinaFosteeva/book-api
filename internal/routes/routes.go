// routes.go
package routes

import (
    "book-api/internal/db"
    "book-api/internal/handlers"
    "github.com/gorilla/mux"
    "log"
)

func SetupRoutes() *mux.Router {
    client, err := db.Connect()
    if err != nil {
        log.Fatal(err)
    }
    collection := client.Database("library").Collection("books")

    app := handlers.NewApp(collection)

    router := mux.NewRouter()
    router.HandleFunc("/books", app.CreateBook).Methods("POST")
    router.HandleFunc("/books", app.GetBooks).Methods("GET")

    return router
}
