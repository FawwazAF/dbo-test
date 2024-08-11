package login

import (
	"context"

	"github.com/dbo-test/internal/model"
	"github.com/dbo-test/pkg/jwt"
)

type repositoryProvider interface {
	GetCustomerByUsername(ctx context.Context, username string) (*model.Customer, error)
}

type jwtItfProvider interface {
	GenerateJWTToken(param jwt.JWTTokenParameter) (string, error)
}

type loginController struct {
	repo repositoryProvider
	jwt  jwtItfProvider
}

func NewLogin(repo repositoryProvider, jwt jwtItfProvider) *loginController {
	return &loginController{repo: repo, jwt: jwt}
}
