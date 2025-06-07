package main

import (
	"context"
	"errors"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	postservice "github.com/JaxonAdams/blog-backend/src/services"
	"github.com/JaxonAdams/blog-backend/src/services/aws/dynamodb"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler(services models.HandlerServices) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		parsedRequest, err := helpers.ParseDeletePostInput(request)
		if err != nil {
			return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
		}

		err = postservice.DeletePost(parsedRequest.ID, services, ctx)
		if err != nil {
			var notFoundErr dynamodb.ErrCodeNotFound
			if errors.As(err, &notFoundErr) {
				return helpers.MakeErrorResponse(404, map[string]string{"message": "Not found"}), nil
			}

			return helpers.MakeErrorResponse(500, map[string]string{"message": err.Error()}), nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 204}, nil
	}
}

func main() {
	handler := createRequestHandler(models.HandlerServices{
		DynamoDBService: dynamodb.New(context.TODO()),
	})
	lambda.Start(handler)
}
