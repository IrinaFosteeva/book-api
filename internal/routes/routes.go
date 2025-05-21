package routes

import (
    "context"
    "log"

    "book-api/internal/handlers"
    "book-api/internal/middleware"
    "book-api/internal/repository"
    "book-api/internal/service"

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

    repo := repository.NewBookRepository(collection)
    svc := service.NewBookService(repo)
    handler := handlers.NewBookHandler(svc)

    router := mux.NewRouter()
    router.Use(middleware.LoggingMiddleware)

    router.HandleFunc("/books", handler.CreateBook).Methods("POST")
    router.HandleFunc("/books", handler.GetBookByID).Queries("id", "{id}").Methods("GET")
    router.HandleFunc("/books", handler.GetBooks).Methods("GET")
    router.HandleFunc("/books", handler.DeleteBook).Queries("id", "{id}").Methods("DELETE")

    return router
}
