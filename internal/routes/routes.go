package routes

import (
    "context"
    "log"
    "book-api/internal/handlers"

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
    router.HandleFunc("/books", app.CreateBook).Methods("POST")
    router.HandleFunc("/books", app.GetBooks).Methods("GET")

    return router
}