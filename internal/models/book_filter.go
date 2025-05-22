package models

type BookFilter struct {
	Title  string
	Author string
	Limit  int64
	Offset int64
}
