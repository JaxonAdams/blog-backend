package models

import "time"

type Post struct {
	ID         string    `json:"id" validate:"required"`
	Title      string    `json:"title" validate:"required"`
	Tags       []string  `json:"tags"`
	HtmlS3Key  string    `json:"html_s3_key"`
	MdS3Key    string    `json:"md_s3_key"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	ModifiedAt time.Time `json:"modified_at" validate:"required"`
}
