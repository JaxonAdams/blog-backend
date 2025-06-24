package main

import (
	"context"
	"strings"

	"github.com/JaxonAdams/blog-backend/src/services/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequestHandler() func(ctx context.Context, request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
		authHeader := request.Headers["authorization"]
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return unauthorized(), nil
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseJWT(tokenString)
		if err != nil {
			return unauthorized(), err
		}

		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: true,
			Context: map[string]any{
				"role": claims.Role,
				"sub":  claims.Subject,
			},
		}, nil
	}
}

func unauthorized() events.APIGatewayV2CustomAuthorizerSimpleResponse {
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
	}
}

func main() {
	handler := createRequestHandler()
	lambda.Start(handler)
}
