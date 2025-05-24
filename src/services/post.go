package postservice

import (
	"fmt"

	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/JaxonAdams/blog-backend/src/services/markdown"
)

func CreatePost(input models.CreatePostInput) (models.Post, error) {
	// Convert the markdown to HTML
	// Store the HTML and Markdown in S3
	// Store metadata in DynamoDB, including S3 key(s)

	html := markdown.MdToHTML([]byte(input.Content))
	fmt.Printf("HTML Content: %s", html)

	return models.Post{
		ID:             "TODO",
		Title:          input.Title,
		Tags:           input.Tags,
		HtmlContentUrl: "TODO",
		MdContentUrl:   "TODO",
	}, nil
}
