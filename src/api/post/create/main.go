package main

import (
	"context"
	"fmt"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler() func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		parsedRequest, err := helpers.ParseCreatePostInput(request)
		if err != nil {
			return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
		}

		fmt.Printf("New Post Title: %s\n", parsedRequest.Title)
		fmt.Printf("New Post Content: %s\n", parsedRequest.Content)
		fmt.Printf("New Post Tags: %v\n", parsedRequest.Tags)

		return helpers.MakeSuccessResponse(201, parsedRequest), nil
	}
}

func main() {
	handler := createRequestHandler()
	lambda.Start(handler)
}
