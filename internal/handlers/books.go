package handlers

import (
    "encoding/json"
    "net/http"
    "book-api/internal/service"
    "github.com/gorilla/mux"
    "github.com/go-playground/validator/v10"
    "book-api/internal/http/apierr"
    "book-api/internal/http/dto"
    "strconv"
    "book-api/internal/models"
)

type BookHandler struct {
    service service.BookService
}

func NewBookHandler(s service.BookService) *BookHandler {
    return &BookHandler{service: s}
}

var validate = validator.New()

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    var input dto.CreateBookInput

    defer r.Body.Close()
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Невалидное тело запроса"))
        return
    }

    if err := validate.Struct(input); err != nil {
        fields := make(map[string]string)
        for _, e := range err.(validator.ValidationErrors) {
            fields[e.Field()] = "некорректное значение"
        }
        apierr.RespondWithError(w, apierr.NewValidationError(fields))
        return
    }

    book, err := h.service.Create(r.Context(), input.Title, input.Author)
    if err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Не удалось создать книгу"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int64(10)
	offset := int64(0)

	if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
		limit = l
	}
	if o, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
		offset = o
	}

	filter := models.BookFilter{
		Title:  title,
		Author: author,
		Limit:  limit,
		Offset: offset,
	}

	books, err := h.service.GetAll(r.Context(), filter)
	if err != nil {
		apierr.RespondWithError(w, apierr.NewInternalError("Не удалось получить книги"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if id == "" {
        apierr.RespondWithError(w, apierr.NewValidationError(map[string]string{
            "id": "обязательный параметр",
        }))
        return
    }

    book, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Книга не найдена или некорректный ID"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    if idStr == "" {
        apierr.RespondWithError(w, apierr.NewValidationError(map[string]string{
            "id": "обязательный параметр",
        }))
        return
    }

    var input dto.UpdateBookInput

    defer r.Body.Close()
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Невалидное тело запроса"))
        return
    }

    if err := validate.Struct(&input); err != nil {
        fields := make(map[string]string)
        for _, e := range err.(validator.ValidationErrors) {
            fields[e.Field()] = "некорректное значение"
        }
        apierr.RespondWithError(w, apierr.NewValidationError(fields))
        return
    }

    book, err := h.service.GetByID(r.Context(), idStr)
    if err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Книга не найдена"))
        return
    }

    if input.Title != nil {
        book.Title = *input.Title
    }
    if input.Author != nil {
        book.Author = *input.Author
    }

    if err := h.service.Update(r.Context(), book); err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Не удалось обновить книгу"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if id == "" {
        apierr.RespondWithError(w, apierr.NewValidationError(map[string]string{
            "id": "обязательный параметр",
        }))
        return
    }

    if err := h.service.DeleteByID(r.Context(), id); err != nil {
        apierr.RespondWithError(w, apierr.NewInternalError("Не удалось удалить книгу"))
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
