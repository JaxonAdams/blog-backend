package models

type Post struct {
	ID        string   `json:"id" validate:"required"`
	Title     string   `json:"title" validate:"required"`
	Tags      []string `json:"tags"`
	HtmlS3Key string   `json:"html_s3_key"`
	MdS3Key   string   `json:"md_s3_key"`
}

type CreatePostInput struct {
	Title   string   `json:"title" validate:"required"`
	Tags    []string `json:"tags"`
	Content string   `json:"content" validate:"required"`
}
