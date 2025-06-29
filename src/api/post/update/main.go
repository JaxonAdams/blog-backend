package main

import (
	"context"
	"errors"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/JaxonAdams/blog-backend/src/services/aws/dynamodb"
	"github.com/JaxonAdams/blog-backend/src/services/aws/s3"
	postservice "github.com/JaxonAdams/blog-backend/src/services/post"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler(services models.HandlerServices) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if !helpers.UserHasAdminRole(request) {
			return helpers.MakeErrorResponse(403, map[string]string{"message": "Forbidden"}), nil
		}

		parsedRequest, err := helpers.ParseUpdatePostInput(request)
		if err != nil {
			return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
		}

		post, err := postservice.UpdatePost(parsedRequest, services, ctx)
		if err != nil {
			var notFoundErr dynamodb.ErrCodeNotFound
			if errors.As(err, &notFoundErr) {
				return helpers.MakeErrorResponse(404, map[string]string{"message": "Not found"}), nil
			}

			var invalidRequestErr postservice.ErrCodeInvalidRequest
			if errors.As(err, &invalidRequestErr) {
				return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
			}

			return helpers.MakeErrorResponse(500, map[string]string{"message": err.Error()}), nil
		}

		return helpers.MakeSuccessResponse(200, map[string]any{"post": post}), nil
	}
}

func main() {
	handler := createRequestHandler(models.HandlerServices{
		S3Service:       s3.New(context.TODO()),
		DynamoDBService: dynamodb.New(context.TODO()),
	})
	lambda.Start(handler)
}
