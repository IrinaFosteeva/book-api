package handlers

import (
    "encoding/json"
    "net/http"
    "book-api/internal/service"
    "log"
    "github.com/gorilla/mux"
)

type BookHandler struct {
    service service.BookService
}

func NewBookHandler(s service.BookService) *BookHandler {
    return &BookHandler{service: s}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title  string `json:"title"`
        Author string `json:"author"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    book, err := h.service.Create(r.Context(), input.Title, input.Author)
    if err != nil {
        http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
    books, err := h.service.GetAll(r.Context())
    if err != nil {
        log.Printf("Failed to fetch books: %v", err)
        http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        http.Error(w, "Missing id parameter", http.StatusBadRequest)
        return
    }

    book, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, "Book not found or invalid ID: "+err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    if idStr == "" {
        http.Error(w, "Missing id parameter", http.StatusBadRequest)
        return
    }

    var input struct {
        Title  *string `json:"title,omitempty"`
        Author *string `json:"author,omitempty"`
    }

    defer r.Body.Close()
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    book, err := h.service.GetByID(r.Context(), idStr)
    if err != nil {
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }

    if input.Title != nil {
        book.Title = *input.Title
    }
    if input.Author != nil {
        book.Author = *input.Author
    }

    err = h.service.Update(r.Context(), book)
    if err != nil {
        http.Error(w, "Failed to update book: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        http.Error(w, "Missing id parameter", http.StatusBadRequest)
        return
    }

    err := h.service.DeleteByID(r.Context(), id)
    if err != nil {
        http.Error(w, "Failed to delete book: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
