package postservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	postmodel "github.com/JaxonAdams/blog-backend/src/models/posts"
	"github.com/JaxonAdams/blog-backend/src/services/markdown"
)

func CreatePost(input models.CreatePostInput, services models.HandlerServices, ctx context.Context) (postmodel.Post, error) {
	// Create a unique ID for the post
	postID := helpers.NewID()

	// Convert the markdown to HTML
	html := markdown.MdToHTML([]byte(input.Content))
	fmt.Printf("HTML Content: %s", html)

	// Store the HTML and Markdown in S3
	htmlS3Key, err := services.S3Service.UploadPostHTML(postID, string(html), ctx)
	if err != nil {
		log.Fatalf("failed to upload md to s3: %v", err)
		return postmodel.Post{}, err
	}

	mdS3Key, err := services.S3Service.UploadPostMd(postID, input.Content, ctx)
	if err != nil {
		log.Fatalf("failed to upload html to s3: %v", err)
		return postmodel.Post{}, err
	}

	post := postmodel.Post{
		ID:         postID,
		Title:      input.Title,
		Tags:       input.Tags,
		HtmlS3Key:  htmlS3Key,
		MdS3Key:    mdS3Key,
		CreatedAt:  time.Now().UnixMilli(),
		ModifiedAt: time.Now().UnixMilli(),
	}

	// Store metadata in DynamoDB, including S3 key
	err = services.DynamoDBService.PutNewPost(post, ctx)
	if err != nil {
		log.Fatalf("failed to store post metadata in dynamo: %v", err)
		return postmodel.Post{}, err
	}

	return post, nil
}

func GetPostByID(id string, services models.HandlerServices, ctx context.Context) (postmodel.Post, error) {
	post, err := services.DynamoDBService.GetPostById(id, ctx)
	if err != nil {
		return postmodel.Post{}, err
	}

	htmlPresignedURL, mdPresignedURL, err := getPresignedUrlsForPost(post, services, ctx)
	if err != nil {
		return postmodel.Post{}, err
	}

	post.HtmlPostUrl = htmlPresignedURL
	post.MdPostUrl = mdPresignedURL

	return post, nil
}

func getPresignedUrlsForPost(post postmodel.Post, services models.HandlerServices, ctx context.Context) (string, string, error) {
	htmlPresignedURL, err := services.S3Service.GetPostHtmlURL(post, ctx)
	if err != nil {
		return "", "", err
	}

	mdPresignedURL, err := services.S3Service.GetPostMdURL(post, ctx)
	if err != nil {
		return "", "", err
	}

	return htmlPresignedURL, mdPresignedURL, nil
}
