package service

import (
    "context"
    "book-api/internal/models"
    "book-api/internal/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService interface {
    Create(ctx context.Context, title, author string) (models.Book, error)
    GetAll(ctx context.Context) ([]models.Book, error)
    GetByID(ctx context.Context, id string) (models.Book, error)
    DeleteByID(ctx context.Context, id string) error
}

type bookService struct {
    repo repository.BookRepository
}

func NewBookService(r repository.BookRepository) BookService {
    return &bookService{repo: r}
}

func (s *bookService) Create(ctx context.Context, title, author string) (models.Book, error) {
    book := models.Book{
        ID:     primitive.NewObjectID(),
        Title:  title,
        Author: author,
    }
    err := s.repo.Create(ctx, book)
    return book, err
}

func (s *bookService) GetAll(ctx context.Context) ([]models.Book, error) {
    return s.repo.GetAll(ctx)
}

func (s *bookService) GetByID(ctx context.Context, id string) (models.Book, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return models.Book{}, err
    }
    return s.repo.GetByID(ctx, objID)
}

func (s *bookService) DeleteByID(ctx context.Context, id string) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    return s.repo.DeleteByID(ctx, objID)
}
