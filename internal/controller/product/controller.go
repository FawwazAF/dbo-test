package product

import (
	"context"

	"github.com/dbo-test/internal/model"
)

type repositoryProvider interface {
	GetAllProduct(ctx context.Context) ([]model.Product, error)
}

type productController struct {
	repo repositoryProvider
}

func NewProduct(repo repositoryProvider) *productController {
	return &productController{repo: repo}
}
