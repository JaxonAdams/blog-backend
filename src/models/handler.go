package models

import "github.com/JaxonAdams/blog-backend/src/services/aws/s3"

type HandlerServices struct {
	S3Service *s3.S3Service
}
