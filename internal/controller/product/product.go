package product

import (
	"context"

	"github.com/dbo-test/internal/model"
)

func (pc *productController) GetAllProduct(ctx context.Context) ([]model.Product, error) {
	return pc.repo.GetAllProduct(ctx)
}
