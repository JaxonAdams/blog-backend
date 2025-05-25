package main

import (
	"context"

	"github.com/JaxonAdams/blog-backend/src/helpers"
	"github.com/JaxonAdams/blog-backend/src/models"
	postservice "github.com/JaxonAdams/blog-backend/src/services"
	"github.com/JaxonAdams/blog-backend/src/services/aws/s3"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler(services models.HandlerServices) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		parsedRequest, err := helpers.ParseCreatePostInput(request)
		if err != nil {
			return helpers.MakeErrorResponse(400, map[string]string{"message": err.Error()}), nil
		}

		createdPost, err := postservice.CreatePost(parsedRequest, services, ctx)
		if err != nil {
			return helpers.MakeErrorResponse(500, map[string]string{"message": err.Error()}), nil
		}

		return helpers.MakeSuccessResponse(201, createdPost), nil
	}
}

func main() {
	handler := createRequestHandler(models.HandlerServices{
		S3Service: s3.New(context.TODO()),
	})
	lambda.Start(handler)
}
