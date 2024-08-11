package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dbo-test/internal/model"
	"github.com/jmoiron/sqlx"
)

func (repo *pgsqlRepository) GetOrderByID(ctx context.Context, orderID, customerID int) (*model.Order, error) {
	query := `
		SELECT 
			id,
			invoice,
			customer_id,
			status,
			created_at,
			updated_at
		FROM
			dbo_trx_order
		WHERE
			id = $1 and customer_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	order := model.Order{}
	if err := repo.pgsql.GetContext(ctx, &order, query, orderID, customerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (repo *pgsqlRepository) GetOrderDetailList(ctx context.Context, orderID, customerID int) ([]model.OrderProduct, error) {
	query := `
		SELECT 
			od.id,
			od.product_id,
			od.quantity,
			od.total_price,
			od.created_at,
			od.updated_at
		FROM
			dbo_dtl_order_detail od LEFT JOIN dbo_trx_order o ON od.order_id = o.id
		WHERE
			o.id = $1 and o.customer_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	orderList := []model.OrderProduct{}
	if err := repo.pgsql.SelectContext(ctx, &orderList, query, orderID, customerID); err != nil {
		return nil, err
	}

	return orderList, nil
}

func (repo *pgsqlRepository) CreateOrder(ctx context.Context, tx SqlTx, customerID int, invoice string) (int, error) {
	query := `
		INSERT INTO dbo_trx_order(
			customer_id,
			invoice,
			status
		) VALUES ($1, $2, $3) RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()
	rows, err := tx.QueryContext(ctx, query, customerID, invoice, model.OrderStatusPending)
	if err != nil {
		return 0, err
	}

	orderID := 0
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&orderID); err != nil {
				return 0, err
			}
		}

		if err := rows.Err(); err != nil {
			return 0, err
		}
	}

	if orderID == 0 {
		return 0, errors.New("got 0 order_id when create order")
	}

	return orderID, nil
}

func (repo *pgsqlRepository) CreateOrderDetail(ctx context.Context, tx SqlTx, orderID int, orderProduct *model.OrderProduct) error {
	query := `
		INSERT INTO dbo_dtl_order_detail(
			order_id,
			product_id,
			quantity,
			total_price
		) VALUES ($1, $2, $3, $4)
	`
	args := []interface{}{
		orderID,
		orderProduct.ProductID,
		orderProduct.Quantity,
		orderProduct.TotalPrice,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) DeleteOrder(ctx context.Context, tx SqlTx, orderID int) error {
	query := `
		DELETE FROM dbo_trx_order WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()
	_, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) DeleteOrderDetail(ctx context.Context, tx SqlTx, orderID int) error {
	query := `
		DELETE FROM dbo_dtl_order_detail WHERE order_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()
	_, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) UpdateOrderDetail(ctx context.Context, tx SqlTx, orderID int, orderProduct *model.OrderProduct) error {
	query := `
		UPDATE dbo_dtl_order_detail SET
		quantity = $1,
		total_price = $2
		WHERE order_id = $3 AND product_id = $4
	`
	args := []interface{}{
		orderProduct.Quantity,
		orderProduct.TotalPrice,
		orderID,
		orderProduct.ProductID,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) SearchOrder(ctx context.Context, customerID int, queryParameters map[string]interface{}) ([]model.Order, error) {
	query, args := repo.buildSearrchOrderQuery(queryParameters)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	fmt.Println(query)

	orders := []model.Order{}
	if err := repo.pgsql.SelectContext(ctx, &orders, query, args...); err != nil {
		return nil, err
	}

	return orders, nil
}

func (repo *pgsqlRepository) buildSearrchOrderQuery(parameters map[string]interface{}) (string, []interface{}) {
	query := `
		SELECT 
			id,
			invoice,
			customer_id,
			status,
			created_at,
			updated_at
		FROM
			dbo_trx_order 
	`
	builder := strings.Builder{}
	args := []interface{}{}
	builder.WriteString(query)

	whereClause := []string{}
	if invoiceValue := parameters["invoice"]; invoiceValue != nil {
		whereClause = append(whereClause, `"invoice" LIKE (?) `)
		args = append(args, "%"+fmt.Sprintf("%v", invoiceValue)+"%")
	}

	if len(whereClause) > 0 {
		builder.WriteString(`WHERE `)
		builder.WriteString(strings.Join(whereClause, "AND "))
	}

	if orderBy := parameters["order_by"]; orderBy != nil {
		builder.WriteString(fmt.Sprintf("ORDER BY id %s ", orderBy))
	}

	if perPage := parameters["per_page"]; perPage != nil {
		builder.WriteString(`LIMIT ? `)
		args = append(args, perPage)
	}

	if page := parameters["page"]; page != nil {
		builder.WriteString(`OFFSET ? `)
		args = append(args, page)
	}

	return sqlx.Rebind(sqlx.DOLLAR, builder.String()), args
}
