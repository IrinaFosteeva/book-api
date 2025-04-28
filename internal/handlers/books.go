package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "book-api/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
    InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
    Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

type App struct {
    collection Collection
}

func NewApp(collection Collection) *App {
    return &App{collection: collection}
}

func (app *App) CreateBook(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title  string `json:"title"`
        Author string `json:"author"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    book := models.Book{
        ID:     primitive.NewObjectID(),
        Title:  input.Title,
        Author: input.Author,
    }

    _, err := app.collection.InsertOne(context.TODO(), book)
    if err != nil {
        http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

func (app *App) GetBooks(w http.ResponseWriter, r *http.Request) {
    cursor, err := app.collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    var books []models.Book
    if err := cursor.All(context.TODO(), &books); err != nil {
        http.Error(w, "Failed to decode books", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(books)
}