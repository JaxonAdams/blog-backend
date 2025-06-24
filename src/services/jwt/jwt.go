package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  username,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return &CustomClaims{}, &ErrCodeInvalidToken{Msg: err.Error()}
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims.Role == "" {
		return &CustomClaims{}, &ErrCodeInvalidToken{Msg: err.Error()}
	}

	return claims, nil
}

type ErrCodeInvalidToken struct {
	Msg string
}

func (e *ErrCodeInvalidToken) Error() string {
	return e.Msg
}
