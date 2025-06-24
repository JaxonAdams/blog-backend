package loginservice

import (
	"context"
	"fmt"

	"github.com/JaxonAdams/blog-backend/src/models"
	"github.com/JaxonAdams/blog-backend/src/services/jwt"
	"golang.org/x/crypto/bcrypt"
)

func LogInAdmin(input models.AdminLoginInput, services models.HandlerServices, ctx context.Context) (string, error) {

	fmt.Printf("Fetching user with input: %+v", input)

	// Fetch the stored password hash from dynamodb
	adminUser, err := services.DynamoDBService.GetAdminUser(input.Username, ctx)
	if err != nil {
		return "", err
	}
	fmt.Printf("Found user with username %s: %+v", input.Username, adminUser)

	// Compare the input password to the fetched hash
	err = bcrypt.CompareHashAndPassword([]byte(adminUser.HashedPW), []byte(input.Password))
	if err != nil {
		return "", ErrCodeUnauthorized{Msg: err.Error()}
	}

	// If correct, generate and return a new JWT
	return jwt.GenerateJWT(adminUser.Username, "admin")
}

type ErrCodeUnauthorized struct {
	Msg string
}

func (e ErrCodeUnauthorized) Error() string {
	return e.Msg
}
