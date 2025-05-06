package routes

import (
    "context"
    "log"
    "book-api/internal/handlers"
    "book-api/internal/middleware"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func SetupRoutes() *mux.Router {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://host.docker.internal:27017"))
    if err != nil {
        log.Fatal(err)
    }
    collection := client.Database("library").Collection("books")

    app := handlers.NewApp(collection)

    router := mux.NewRouter()
    router.Use(middleware.LoggingMiddleware)
    router.HandleFunc("/books", app.CreateBook).Methods("POST")
    router.HandleFunc("/books", app.GetBookByID).Queries("id", "{id}").Methods("GET")
    router.HandleFunc("/books", app.GetBooks).Methods("GET")
    router.HandleFunc("/books", app.DeleteBook).Queries("id", "{id}").Methods("DELETE")

    return router
}