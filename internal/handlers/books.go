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
    FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
    DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
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

func (app *App) GetBookByID(w http.ResponseWriter, r *http.Request) {
    vars := r.URL.Query()
    id := vars.Get("id")
    if id == "" {
        http.Error(w, "Missing id parameter", http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid id format", http.StatusBadRequest)
        return
    }

    var book models.Book
    result := app.collection.FindOne(context.TODO(), bson.M{"_id": objID})
    if err := result.Decode(&book); err == mongo.ErrNoDocuments {
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, "Failed to decode book: "+err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(book)
}

func (app *App) DeleteBook(w http.ResponseWriter, r *http.Request) {
    vars := r.URL.Query()
    id := vars.Get("id")
    if id == "" {
        http.Error(w, "Missing id parameter", http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid id format", http.StatusBadRequest)
        return
    }

    result, err := app.collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        http.Error(w, "Failed to delete book: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if result.DeletedCount == 0 {
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}