package models

type CreatePostInput struct {
	Title   string   `json:"title" validate:"required"`
	Tags    []string `json:"tags"`
	Content string   `json:"content" validate:"required"`
}

type GetPostByIdInput struct {
	ID string `json:"string" validate:"required"`
}
