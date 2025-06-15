package main

import (
	"context"
	"fmt"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
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

		fmt.Printf(
			"Received login request: {username: '%s', password: '%s'}",
			parsedRequest.Username,
			parsedRequest.Password,
		)

		loginservice.LogInAdmin(parsedRequest, services, ctx)

		return helpers.MakeSuccessResponse(200, map[string]any{"message": "Hello, world!"}), nil
	}
}

func main() {
	handler := createRequestHandler(models.HandlerServices{})
	lambda.Start(handler)
}
