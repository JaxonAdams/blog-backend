package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/aws/aws-lambda-go/events"
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
