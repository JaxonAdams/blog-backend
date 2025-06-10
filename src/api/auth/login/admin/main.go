package main

import (
	"context"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler() func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return helpers.MakeSuccessResponse(200, map[string]any{"message": "Hello, world!"}), nil
	}
}

func main() {
	handler := createRequestHandler()
	lambda.Start(handler)
}
