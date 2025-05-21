package main

import (
	"context"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/services/markdown"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler() func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		parsedRequest, err := helpers.ParseCreatePostInput(request)
		if err != nil {
			return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
		}

		html := markdown.MdToHTML([]byte(parsedRequest.Content))

		response := map[string]string{
			"html": string(html),
		}
		return helpers.MakeSuccessResponse(201, response), nil
	}
}

func main() {
	handler := createRequestHandler()
	lambda.Start(handler)
}
