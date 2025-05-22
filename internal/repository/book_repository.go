package repository

import (
    "context"
    "book-api/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookRepository interface {
    Create(ctx context.Context, book models.Book) error
    GetAll(ctx context.Context, filter models.BookFilter) ([]models.Book, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (models.Book, error)
    DeleteByID(ctx context.Context, id primitive.ObjectID) error
	Update(ctx context.Context, book models.Book) error
}

type bookRepo struct {
    collection *mongo.Collection
}

func NewBookRepository(col *mongo.Collection) BookRepository {
    return &bookRepo{collection: col}
}

func (r *bookRepo) Create(ctx context.Context, book models.Book) error {
    _, err := r.collection.InsertOne(ctx, book)
    return err
}

func (r *bookRepo) Update(ctx context.Context, book models.Book) error {
    filter := bson.M{"_id": book.ID}
    update := bson.M{
        "$set": bson.M{
            "title":  book.Title,
            "author": book.Author,
        },
    }

    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return mongo.ErrNoDocuments
    }
    return nil
}

func (r *bookRepo) GetAll(ctx context.Context, filter models.BookFilter) ([]models.Book, error) {
	bsonFilter := bson.M{}
	if filter.Title != "" {
		bsonFilter["title"] = bson.M{"$regex": filter.Title, "$options": "i"}
	}
	if filter.Author != "" {
		bsonFilter["author"] = bson.M{"$regex": filter.Author, "$options": "i"}
	}

	findOptions := options.Find()
	if filter.Limit > 0 {
		findOptions.SetLimit(filter.Limit)
	}
	findOptions.SetSkip(filter.Offset)

	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (r *bookRepo) GetByID(ctx context.Context, id primitive.ObjectID) (models.Book, error) {
    var book models.Book
    result := r.collection.FindOne(ctx, bson.M{"_id": id})
    err := result.Decode(&book)
    return book, err
}

func (r *bookRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
    result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return mongo.ErrNoDocuments
    }
    return nil
}
