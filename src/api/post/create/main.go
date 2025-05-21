package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler() func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
		fmt.Printf("Body size = %d.\n", len(request.Body))

		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("\t%s: %s\n", key, value)
		}

		return events.APIGatewayProxyResponse{
			Body:       request.Body,
			StatusCode: 200,
		}, nil
	}
}

func main() {
	handler := createRequestHandler()
	lambda.Start(handler)
}
