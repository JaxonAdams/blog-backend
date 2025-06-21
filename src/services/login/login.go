package loginservice

import (
	"context"
	"fmt"

	"github.com/JaxonAdams/blog-backend/src/models"
	usermodel "github.com/JaxonAdams/blog-backend/src/models/users"
)

func LogInAdmin(input models.AdminLoginInput, services models.HandlerServices, ctx context.Context) (usermodel.AdminUser, error) {

	fmt.Printf("Fetching user with input: %+v", input)

	// Fetch the stored password hash from dynamodb
	adminUser, err := services.DynamoDBService.GetAdminUser(input.Username, ctx)
	if err != nil {
		return usermodel.AdminUser{}, err
	}
	fmt.Printf("Found user with username %s: %+v", input.Username, adminUser)

	// Compare the input password to the fetched hash
	// If correct, generate and return a new JWT

	return adminUser, nil
}
