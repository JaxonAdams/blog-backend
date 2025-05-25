package postservice

import (
	"context"
	"fmt"
	"log"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/JaxonAdams/blog-backend/src/services/markdown"
)

func CreatePost(input models.CreatePostInput, services models.HandlerServices, ctx context.Context) (models.Post, error) {
	// Create a unique ID for the post
	postID := helpers.NewID()

	// Convert the markdown to HTML
	html := markdown.MdToHTML([]byte(input.Content))
	fmt.Printf("HTML Content: %s", html)

	// Store the HTML and Markdown in S3
	htmlS3Key, err := services.S3Service.UploadPostMd(postID, input.Content, ctx)
	if err != nil {
		log.Fatalf("failed to upload md to s3: %v", err)
		return models.Post{}, err
	}

	mdS3Key, err := services.S3Service.UploadPostHTML(postID, string(html), ctx)
	if err != nil {
		log.Fatalf("failed to upload html to s3: %v", err)
		return models.Post{}, err
	}

	// TODO: Store metadata in DynamoDB, including S3 key

	return models.Post{
		ID:        postID,
		Title:     input.Title,
		Tags:      input.Tags,
		HtmlS3Key: htmlS3Key,
		MdS3Key:   mdS3Key,
	}, nil
}
