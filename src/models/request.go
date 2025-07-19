package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type CreatePostInput struct {
	Title   string   `json:"title" validate:"required"`
	Summary string   `json:"summary" validate:"required"`
	Tags    []string `json:"tags" validate:"required"`
	Content string   `json:"content" validate:"required"`
}

type GetPostByIdInput struct {
	ID string `json:"string" validate:"required"`
}

type GetPostsInput struct {
	PageSize int
	StartKey map[string]types.AttributeValue
}

type UpdatePostInput struct {
	GetPostByIdInput
	Title   *string   `json:"title" validate:"required"`
	Summary *string   `json:"summary" validate:"required"`
	Tags    *[]string `json:"tags" validate:"required"`
	Content *string   `json:"content" validate:"required"`
}

type DeletePostInput struct {
	GetPostByIdInput
}

type AdminLoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
