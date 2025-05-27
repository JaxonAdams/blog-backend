package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func ParseCreatePostInput(request events.APIGatewayProxyRequest) (models.CreatePostInput, error) {
	var input models.CreatePostInput

	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return models.CreatePostInput{}, err
	}

	if input.Title == "" || input.Content == "" {
		return models.CreatePostInput{}, fmt.Errorf("fields title and content are required")
	}

	return input, nil
}

func ParseGetPostByIdInput(request events.APIGatewayProxyRequest) (models.GetPostByIdInput, error) {
	pathParams := request.PathParameters

	id, exists := pathParams["post_id"]
	if !exists {
		return models.GetPostByIdInput{}, fmt.Errorf("post_id path param is required")
	}

	return models.GetPostByIdInput{
		ID: id,
	}, nil
}

func ParseGetPostsInput(request events.APIGatewayProxyRequest) (models.GetPostsInput, error) {
	var startKey map[string]types.AttributeValue
	pageSize, _ := strconv.Atoi(os.Getenv("DEFAULT_PAGE_SIZE"))

	queryStringParams := request.QueryStringParameters

	if v, exists := queryStringParams["pageSize"]; exists {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	if v, exists := queryStringParams["startKey"]; exists && v != "" {
		sk, err := decodeStartKey(v)
		if err != nil {
			return models.GetPostsInput{}, err
		}
		startKey = sk
	}

	return models.GetPostsInput{
		PageSize: pageSize,
		StartKey: startKey,
	}, nil
}

func MakeSuccessResponse(statusCode int, data any) events.APIGatewayProxyResponse {
	response := map[string]any{
		"data": data,
	}

	body, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func MakeErrorResponse(statusCode int, errorInfo any) events.APIGatewayProxyResponse {
	response := map[string]any{
		"error": errorInfo,
	}

	body, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func decodeStartKey(encoded string) (map[string]types.AttributeValue, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("startKey is not valid base64: %w", err)
	}

	var raw map[string]map[string]string
	if err := json.Unmarshal(decoded, &raw); err != nil {
		return nil, fmt.Errorf("invalid startKey JSON: %w", err)
	}

	startKey := make(map[string]types.AttributeValue)
	for k, val := range raw {
		if s, ok := val["S"]; ok {
			startKey[k] = &types.AttributeValueMemberS{Value: s}
		} else if n, ok := val["N"]; ok {
			startKey[k] = &types.AttributeValueMemberN{Value: n}
		} else {
			return nil, fmt.Errorf("unsupported attribute type for key: %s", k)
		}
	}
	return startKey, nil
}
