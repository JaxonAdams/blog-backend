package main

import (
	"context"
	"errors"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/JaxonAdams/blog-backend/src/services/aws/dynamodb"
	loginservice "github.com/JaxonAdams/blog-backend/src/services/login"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler(services models.HandlerServices) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		parsedRequest, err := helpers.ParseAdminLoginInput(request)
		if err != nil {
			return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
		}

		token, err := loginservice.LogInAdmin(parsedRequest, services, ctx)
		if err != nil {
			var notFoundErr dynamodb.ErrCodeNotFound
			var unauthorizedError loginservice.ErrCodeUnauthorized
			if errors.As(err, &notFoundErr) || errors.As(err, &unauthorizedError) {
				return helpers.MakeErrorResponse(401, map[string]string{"message": "Unauthorized"}), nil
			}
			return helpers.MakeErrorResponse(500, map[string]string{"message": err.Error()}), nil
		}

		return helpers.MakeSuccessResponse(200, map[string]any{"token": token}), nil
	}
}

func main() {
	services := models.HandlerServices{}
	services.DynamoDBService = dynamodb.New(context.TODO())

	handler := createRequestHandler(services)
	lambda.Start(handler)
}
