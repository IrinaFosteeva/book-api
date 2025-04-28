package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
    ID     primitive.ObjectID `json:"id" bson:"_id"`
    Title  string             `json:"title" bson:"title"`
    Author string             `json:"author" bson:"author"`
}