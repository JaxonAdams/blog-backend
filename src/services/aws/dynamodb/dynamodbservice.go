package dynamodb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	postmodel "github.com/JaxonAdams/blog-backend/src/models/posts"
	usermodel "github.com/JaxonAdams/blog-backend/src/models/users"
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

func (d DynamoDBService) UpsertPost(post postmodel.Post, ctx context.Context) error {
	table := os.Getenv("POST_METADATA_TABLE_NAME")
	item := post.DynamoFormat()

	return d.putItem(table, item, ctx)
}

func (d DynamoDBService) DeletePost(postId string, createdAt int, ctx context.Context) error {
	table := os.Getenv("POST_METADATA_TABLE_NAME")
	return d.deleteItem(table, postId, createdAt, ctx)
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

	var nextStartKey string
	if result.LastEvaluatedKey != nil {
		// Convert raw AttributeValues to an intermediate JSON-safe map
		jsonFriendlyKey := make(map[string]map[string]string)
		for k, v := range result.LastEvaluatedKey {
			switch attr := v.(type) {
			case *types.AttributeValueMemberS:
				jsonFriendlyKey[k] = map[string]string{"S": attr.Value}
			case *types.AttributeValueMemberN:
				jsonFriendlyKey[k] = map[string]string{"N": attr.Value}
			default:
				continue
			}
		}

		// Encode in JSON, then base64
		startKeyJson, _ := json.Marshal(jsonFriendlyKey)
		encodedStartKey := base64.StdEncoding.EncodeToString(startKeyJson)
		nextStartKey = encodedStartKey
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

func (d DynamoDBService) GetAdminUser(username string, ctx context.Context) (usermodel.AdminUser, error) {
	table := os.Getenv("AUTH_TABLE_NAME")

	input := &dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: username},
		},
		Limit: aws.Int32(1),
	}

	result, err := d.client.Query(ctx, input)
	if err != nil {
		return usermodel.AdminUser{}, nil
	}

	if len(result.Items) == 0 {
		fmt.Printf("No users found with username %s", username)
		return usermodel.AdminUser{}, ErrCodeNotFound{Msg: fmt.Sprintf("no user found with username %s", username)}
	}

	var user usermodel.AdminUser
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return usermodel.AdminUser{}, err
	}

	return user, nil
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

func (d DynamoDBService) deleteItem(tableName, itemId string, itemCreatedAt int, ctx context.Context) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: itemId},
			"createdAt": &types.AttributeValueMemberN{Value: strconv.Itoa(itemCreatedAt)},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	}

	_, err := d.client.DeleteItem(ctx, input)
	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if ok := errors.As(err, &cce); ok {
			return ErrCodeNotFound{Msg: fmt.Sprintf("no post found with id %s", itemId)}
		}
		return fmt.Errorf("failed to delete post: %w", err)
	}

	log.Printf("Post with ID %s successfully deleted.", itemId)
	return nil
}

type ErrCodeNotFound struct {
	Msg string
}

func (e ErrCodeNotFound) Error() string {
	return e.Msg
}
