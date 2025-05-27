package dynamodb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	postmodel "github.com/JaxonAdams/blog-backend/src/models/posts"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func (d DynamoDBService) GetAllPosts(pageSize int32, startKey map[string]types.AttributeValue, ctx context.Context) ([]postmodel.Post, string, error) {
	table := os.Getenv("POST_METADATA_TABLE_NAME")
	input := &dynamodb.ScanInput{
		TableName:         &table,
		Limit:             &pageSize,
		ExclusiveStartKey: startKey,
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return []postmodel.Post{}, "", err
	}

	var posts []postmodel.Post
	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return posts, "", err
	}

	// TODO: fix start key issues
	var nextStartKey string
	if result.LastEvaluatedKey != nil {
		marshaledKey, err := attributevalue.MarshalMap(result.LastEvaluatedKey)
		if err == nil {
			startKeyJson, _ := json.Marshal(marshaledKey)
			encodedStartKey := base64.StdEncoding.EncodeToString(startKeyJson)
			nextStartKey = encodedStartKey
		}
	}

	return posts, nextStartKey, nil
}

func (d DynamoDBService) GetPostById(id string, ctx context.Context) (postmodel.Post, error) {
	table := os.Getenv("POST_METADATA_TABLE_NAME")

	input := &dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: id},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	}

	result, err := d.client.Query(ctx, input)
	if err != nil {
		return postmodel.Post{}, err
	}

	if len(result.Items) == 0 {
		fmt.Printf("No posts found with ID %s", id)
		return postmodel.Post{}, ErrCodeNotFound{Msg: fmt.Sprintf("no post found with id %s", id)}
	}

	var post postmodel.Post
	err = attributevalue.UnmarshalMap(result.Items[0], &post)
	if err != nil {
		return postmodel.Post{}, err
	}

	return post, nil
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

type ErrCodeNotFound struct {
	Msg string
}

func (e ErrCodeNotFound) Error() string {
	return e.Msg
}
