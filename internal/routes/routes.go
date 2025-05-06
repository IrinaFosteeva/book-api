package routes

import (
    "context"
    "log"
    "net/http"
    "time"
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
    router.Use(loggingMiddleware)
    router.HandleFunc("/books", app.CreateBook).Methods("POST")
    router.HandleFunc("/books", app.GetBookByID).Queries("id", "{id}").Methods("GET")
    router.HandleFunc("/books", app.GetBooks).Methods("GET")
    router.HandleFunc("/books", app.DeleteBook).Queries("id", "{id}").Methods("DELETE")

    return router
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.String())
        next.ServeHTTP(w, r)
        log.Printf("Completed %s %s in %v", r.Method, r.URL.String(), time.Since(start))
    })
}