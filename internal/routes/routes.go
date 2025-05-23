package routes

import (
    "log"
    "book-api/internal/handlers"
    "book-api/internal/middleware"
    "book-api/internal/repository"
    "book-api/internal/service"
    "book-api/internal/db"
    "github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    client, err := db.Connect()
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
    router.HandleFunc("/books/{id}", handler.GetBookByID).Methods("GET")
    router.HandleFunc("/books", handler.GetBooks).Methods("GET")
    router.HandleFunc("/books/{id}", handler.DeleteBook).Methods("DELETE")
    router.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PATCH")


    return router
}
