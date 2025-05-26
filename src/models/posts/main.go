package postmodel

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Post struct {
	ID         string   `json:"id" validate:"required"`
	Title      string   `json:"title" validate:"required"`
	Tags       []string `json:"tags"`
	HtmlS3Key  string   `json:"html_s3_key" dynamodbav:"html_s3_key"`
	MdS3Key    string   `json:"md_s3_key" dynamodbav:"md_s3_key"`
	CreatedAt  int64    `json:"created_at" validate:"required"`
	ModifiedAt int64    `json:"modified_at" validate:"required"`
}

func (p Post) DynamoFormat() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: p.ID},
		"title":       &types.AttributeValueMemberS{Value: p.Title},
		"tags":        &types.AttributeValueMemberSS{Value: p.Tags},
		"html_s3_key": &types.AttributeValueMemberS{Value: p.HtmlS3Key},
		"md_s3_key":   &types.AttributeValueMemberS{Value: p.MdS3Key},
		"createdAt":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", p.CreatedAt)},
		"modifiedAt":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", p.ModifiedAt)},
	}
}
