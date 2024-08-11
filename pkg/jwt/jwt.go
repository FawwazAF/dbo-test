package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	jwtSecrets string
}

func NewJWT(jwtToken string) *JWT {
	return &JWT{
		jwtSecrets: jwtToken,
	}
}

func (j *JWT) GetJWTSecret() []byte {
	return []byte(j.jwtSecrets)
}

func (j *JWT) ParseClientToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

type JWTTokenParameter struct {
	ID             int
	Username       string
	ExpirationDate time.Time
}

func (j *JWT) GenerateJWTToken(param JWTTokenParameter) (string, error) {
	// Create a new token with user claims.
	claims := jwt.MapClaims{
		"id":       param.ID,
		"username": param.Username,
		"exp":      param.ExpirationDate.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.GetJWTSecret())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
