package models

type Post struct {
	ID             string   `json:"id" validate:"required"`
	Title          string   `json:"title" validate:"required"`
	Tags           []string `json:"tags"`
	HtmlContentUrl string   `json:"html_content_url"`
	MdContentUrl   string   `json:"md_content_url"`
}

type CreatePostInput struct {
	Title   string   `json:"title" validate:"required"`
	Tags    []string `json:"tags"`
	Content string   `json:"content" validate:"required"`
}
