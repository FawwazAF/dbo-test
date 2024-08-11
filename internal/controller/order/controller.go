package order

import (
	"context"

	"github.com/dbo-test/internal/model"
	"github.com/dbo-test/internal/repository/pgsql"
)

type repository interface {
	// Transaction
	BeginTx(ctx context.Context) (pgsql.SqlTx, error)

	// Order Repository
	GetOrderByID(ctx context.Context, orderID, customerID int) (*model.Order, error)
	GetOrderDetailList(ctx context.Context, orderID, customerID int) ([]model.OrderProduct, error)
	CreateOrder(ctx context.Context, tx pgsql.SqlTx, customerID int, invoice string) (int, error)
	CreateOrderDetail(ctx context.Context, tx pgsql.SqlTx, orderID int, orderProduct *model.OrderProduct) error
	DeleteOrder(ctx context.Context, tx pgsql.SqlTx, orderID int) error
	DeleteOrderDetail(ctx context.Context, tx pgsql.SqlTx, orderID int) error
	UpdateOrderDetail(ctx context.Context, tx pgsql.SqlTx, orderID int, orderProduct *model.OrderProduct) error
	SearchOrder(ctx context.Context, customerID int, queryParameters map[string]interface{}) ([]model.Order, error)

	// Customer Repository
	GetCustomerByID(ctx context.Context, id int) (*model.Customer, error)

	// Product repository
	GetProductDetailByID(ctx context.Context, productID int) (*model.Product, error)
	UpdateProductStock(ctx context.Context, tx pgsql.SqlTx, productID, remainingStock int) error
}

type orderController struct {
	repo repository
}

func NewOrder(repo repository) *orderController {
	return &orderController{repo: repo}
}
