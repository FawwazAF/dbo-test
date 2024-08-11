package login

import (
	"context"
	"errors"
	"time"

	"github.com/dbo-test/pkg/hash"
	"github.com/dbo-test/pkg/jwt"
)

func (lc *loginController) Login(ctx context.Context, username, password string) (string, error) {
	customer, err := lc.repo.GetCustomerByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if customer.ID == 0 {
		return "", errors.New("username not found")
	}

	if !hash.CheckPasswordHash(password, string(customer.GetHashedPassword())) {
		return "", errors.New("invalid password")
	}

	token, err := lc.jwt.GenerateJWTToken(jwt.JWTTokenParameter{
		ID:             customer.ID,
		Username:       username,
		ExpirationDate: time.Now().Add(time.Duration(60 * time.Minute)),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
