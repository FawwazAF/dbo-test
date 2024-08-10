package customer

import (
	"context"

	"github.com/dbo-test/internal/model"
)

type customerRepository interface {
	GetCustomerByID(ctx context.Context, id int) (*model.Customer, error)
	GetCustomerByUsername(ctx context.Context, username string) (*model.Customer, error)
	AddCustomer(ctx context.Context, req *model.Customer) error
	UpdateCustomer(ctx context.Context, req *model.Customer) error
	DeleteCustomer(ctx context.Context, id int) error
}

type customerController struct {
	customer customerRepository
}

func NewCustomer(customer customerRepository) *customerController {
	return &customerController{customer: customer}
}
