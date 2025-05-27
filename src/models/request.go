package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type CreatePostInput struct {
	Title   string   `json:"title" validate:"required"`
	Tags    []string `json:"tags"`
	Content string   `json:"content" validate:"required"`
}

type GetPostByIdInput struct {
	ID string `json:"string" validate:"required"`
}

type GetPostsInput struct {
	PageSize int
	StartKey map[string]types.AttributeValue
}
