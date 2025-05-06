package tests

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "book-api/internal/handlers"
    "book-api/internal/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct{}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
    return &mongo.InsertOneResult{}, nil
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
    cursor, err := mongo.NewCursorFromDocuments(nil, nil, nil)
    if err != nil {
        return nil, err
    }
    return cursor, nil
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
    book := models.Book{
        ID:     primitive.NewObjectID(),
        Title:  "Test Book",
        Author: "Test Author",
    }
    return mongo.NewSingleResultFromDocument(book, nil, nil)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
    return &mongo.DeleteResult{DeletedCount: 1}, nil
}

func TestCreateBook(t *testing.T) {
    mockCollection := &MockCollection{}
    app := handlers.NewApp(mockCollection)

    input := struct {
        Title  string `json:"title"`
        Author string `json:"author"`
    }{Title: "Test Book", Author: "Test Author"}
    body, _ := json.Marshal(input)

    req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    app.CreateBook(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("Expected status %v, got %v", http.StatusCreated, status)
    }

    var response models.Book
    if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
        t.Errorf("Failed to decode response: %v", err)
    }
    if response.Title != input.Title || response.Author != input.Author {
        t.Errorf("Expected book {Title:%q, Author:%q}, got %+v", input.Title, input.Author, response)
    }
    if response.ID.IsZero() {
        t.Errorf("Expected non-zero ObjectID, got zero")
    }
}

func TestGetBookByID(t *testing.T) {
    mockCollection := &MockCollection{}
    app := handlers.NewApp(mockCollection)

    req, _ := http.NewRequest("GET", "/books?id=507f1f77bcf86cd799439011", nil)
    rr := httptest.NewRecorder()

    app.GetBookByID(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Expected status %v, got %v", http.StatusOK, status)
    }

    var response models.Book
    if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
        t.Errorf("Failed to decode response: %v", err)
    }
    if response.Title != "Test Book" || response.Author != "Test Author" {
        t.Errorf("Expected book {Title:%q, Author:%q}, got %+v", "Test Book", "Test Author", response)
    }
}

func TestDeleteBook(t *testing.T) {
    mockCollection := &MockCollection{}
    app := handlers.NewApp(mockCollection)

    req, _ := http.NewRequest("DELETE", "/books?id=507f1f77bcf86cd799439011", nil)
    rr := httptest.NewRecorder()

    app.DeleteBook(rr, req)

    if status := rr.Code; status != http.StatusNoContent {
        t.Errorf("Expected status %v, got %v", http.StatusNoContent, status)
    }
}