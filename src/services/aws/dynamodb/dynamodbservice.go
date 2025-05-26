package dynamodb

import (
	"context"
	"log"
	"os"

	postmodel "github.com/JaxonAdams/blog-backend/src/models/posts"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBService struct {
	client *dynamodb.Client
}

func New(ctx context.Context) *DynamoDBService {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBService{
		client: client,
	}
}

func (d DynamoDBService) PutNewPost(post postmodel.Post, ctx context.Context) error {
	table := os.Getenv("POST_METADATA_TABLE_NAME")
	item := post.DynamoFormat()

	return d.putItem(table, item, ctx)
}

func (d DynamoDBService) putItem(tableName string, item map[string]types.AttributeValue, ctx context.Context) error {
	input := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	}

	_, err := d.client.PutItem(ctx, input)
	if err != nil {
		return err
	}

	log.Println("Post successfully stored in DynamoDB")
	return nil
}
