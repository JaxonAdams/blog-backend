package main

import (
	"context"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/JaxonAdams/blog-backend/src/services/aws/dynamodb"
	"github.com/JaxonAdams/blog-backend/src/services/aws/s3"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler(services models.HandlerServices) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return helpers.MakeSuccessResponse(200, map[string]any{"message": "Hello, world!"}), nil
	}
}

func main() {
	handler := createRequestHandler(models.HandlerServices{
		S3Service:       s3.New(context.TODO()),
		DynamoDBService: dynamodb.New(context.TODO()),
	})
	lambda.Start(handler)
}
