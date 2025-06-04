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

func (d DynamoDBService) UpdatePost(post *postmodel.PartialPostUpdate, ctx context.Context) error {
	table := os.Getenv("POST_METADATA_TABLE_NAME")
	input, err := buildUpdateInput(table, post)
	if err != nil {
		return err
	}

	return d.updateItem(input, ctx)
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

func (d DynamoDBService) updateItem(input *dynamodb.UpdateItemInput, ctx context.Context) error {
	_, err := d.client.UpdateItem(ctx, input)
	if err != nil {
		return err
	}

	log.Println("Post successfully stored in DynamoDB")
	return nil
}

func buildUpdateInput(tableName string, post *postmodel.PartialPostUpdate) (*dynamodb.UpdateItemInput, error) {
	if post.ID == "" {
		return nil, fmt.Errorf("missing ID")
	}

	exprAttrValues := map[string]types.AttributeValue{}
	exprAttrNames := map[string]string{}
	updateExprParts := []string{}

	// Optional fields
	hasOptionalUpdates := addOptionalFieldUpdates(post, &updateExprParts, exprAttrValues, exprAttrNames)

	// Always-updated fields
	addAlwaysUpdatedFields(post, &updateExprParts, exprAttrValues, exprAttrNames)

	// Combine update expressions
	updateExpr := "SET " + joinUpdateExpr(updateExprParts)

	if !hasOptionalUpdates {
		log.Println("No optional fields provided; updating only default fields.")
	}

	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: post.ID}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrValues,
		ExpressionAttributeNames:  exprAttrNames,
		ReturnValues:              types.ReturnValueAllNew,
	}, nil
}

func addOptionalFieldUpdates(post *postmodel.PartialPostUpdate, updateExprParts *[]string, exprAttrValues map[string]types.AttributeValue, exprAttrNames map[string]string) bool {
	hasUpdate := false

	if post.Title != nil {
		*updateExprParts = append(*updateExprParts, "#title = :title")
		exprAttrValues[":title"] = &types.AttributeValueMemberS{Value: *post.Title}
		exprAttrNames["#title"] = "title"
		hasUpdate = true
	}

	if post.Tags != nil {
		tagList := make([]types.AttributeValue, len(*post.Tags))
		for i, t := range *post.Tags {
			tagList[i] = &types.AttributeValueMemberS{Value: t}
		}
		*updateExprParts = append(*updateExprParts, "#tags = :tags")
		exprAttrValues[":tags"] = &types.AttributeValueMemberL{Value: tagList}
		exprAttrNames["#tags"] = "tags"
		hasUpdate = true
	}

	return hasUpdate
}

func addAlwaysUpdatedFields(post *postmodel.PartialPostUpdate, updateExprParts *[]string, exprAttrValues map[string]types.AttributeValue, exprAttrNames map[string]string) {
	*updateExprParts = append(*updateExprParts, "#htmlPostUrl = :htmlPostUrl")
	exprAttrValues[":htmlPostUrl"] = &types.AttributeValueMemberS{Value: post.HtmlPostUrl}
	exprAttrNames["#htmlPostUrl"] = "html_post_url"

	*updateExprParts = append(*updateExprParts, "#mdPostUrl = :mdPostUrl")
	exprAttrValues[":mdPostUrl"] = &types.AttributeValueMemberS{Value: post.MdPostUrl}
	exprAttrNames["#mdPostUrl"] = "md_post_url"

	*updateExprParts = append(*updateExprParts, "#htmlS3Key = :htmlS3Key")
	exprAttrValues[":htmlS3Key"] = &types.AttributeValueMemberS{Value: post.HtmlS3Key}
	exprAttrNames["#htmlS3Key"] = "html_s3_key"

	*updateExprParts = append(*updateExprParts, "#mdS3Key = :mdS3Key")
	exprAttrValues[":mdS3Key"] = &types.AttributeValueMemberS{Value: post.MdS3Key}
	exprAttrNames["#mdS3Key"] = "md_s3_key"

	*updateExprParts = append(*updateExprParts, "#createdAt = :createdAt")
	exprAttrValues[":createdAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", post.CreatedAt)}
	exprAttrNames["#createdAt"] = "created_at"

	*updateExprParts = append(*updateExprParts, "#modifiedAt = :modifiedAt")
	exprAttrValues[":modifiedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", post.ModifiedAt)}
	exprAttrNames["#modifiedAt"] = "modified_at"
}

func joinUpdateExpr(parts []string) string {
	return joinStrings(parts, ", ")
}

func joinStrings(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for _, s := range slice[1:] {
		result += sep + s
	}
	return result
}

type ErrCodeNotFound struct {
	Msg string
}

func (e ErrCodeNotFound) Error() string {
	return e.Msg
}
