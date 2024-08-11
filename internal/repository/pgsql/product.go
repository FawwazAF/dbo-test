package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dbo-test/internal/model"
)

func (repo *pgsqlRepository) GetProductDetailByID(ctx context.Context, productID int) (*model.Product, error) {
	query := `
		SELECT 
			id,
			price,
			stock
		FROM
			dbo_mst_product
		WHERE
			id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	product := model.Product{}
	if err := repo.pgsql.GetContext(ctx, &product, query, productID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (repo *pgsqlRepository) UpdateProductStock(ctx context.Context, tx SqlTx, productID, remainingStock int) error {
	query := `
		UPDATE dbo_mst_product SET 
		stock = $1
		WHERE id = $2
	`
	args := []interface{}{
		remainingStock,
		productID,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
