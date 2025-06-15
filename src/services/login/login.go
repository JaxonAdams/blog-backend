package loginservice

import (
	"context"

	"github.com/JaxonAdams/blog-backend/src/models"
)

func LogInAdmin(input models.AdminLoginInput, services models.HandlerServices, ctx context.Context) error {

	// Fetch the stored password hash from dynamodb
	// Compare the input password to the fetched hash
	// If correct, generate and return a new JWT

	return nil
}
