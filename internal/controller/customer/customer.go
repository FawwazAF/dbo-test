package customer

import (
	"context"
	"errors"

	"github.com/dbo-test/internal/model"
	"github.com/dbo-test/pkg/hash"
)

func (cc *customerController) GetCustomerByID(ctx context.Context, id int) (*model.Customer, error) {
	return cc.customer.GetCustomerByID(ctx, id)
}

func (cc *customerController) AddCustomer(ctx context.Context, req *model.Customer) error {
	// get existing username
	customer, err := cc.customer.GetCustomerByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	// if username exist, return error
	if customer.ID != 0 {
		return errors.New("username already exist")
	}

	// encrypt password
	hashedPassword, err := hash.HashPassword(req.GetPassword())
	if err != nil {
		return err
	}
	req.SetHashedPassword(hashedPassword)

	// set customer status active
	req.Status = 1

	return cc.customer.AddCustomer(ctx, req)
}

func (cc *customerController) UpdateCustomer(ctx context.Context, req *model.Customer) error {
	return cc.customer.UpdateCustomer(ctx, req)
}

func (cc *customerController) DeleteCustomer(ctx context.Context, id int) error {
	return cc.customer.DeleteCustomer(ctx, id)
}
