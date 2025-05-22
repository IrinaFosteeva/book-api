package dto

type CreateBookInput struct {
    Title  string `json:"title" validate:"required,min=2,max=50"`
    Author string `json:"author" validate:"required,min=2,max=50"`
}

type UpdateBookInput struct {
    Title  *string `json:"title,omitempty" validate:"omitempty,min=2,max=50"`
    Author *string `json:"author,omitempty" validate:"omitempty,min=2,max=50"`
}
