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
	hashedPassword, err := hash.HashPassword(req.Password)
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

func (cc *customerController) SearchCustomer(ctx context.Context, query map[string]interface{}) ([]model.Customer, bool, error) {
	// adjust pagination
	var (
		page    int
		perPage int
	)
	pageRaw := query["page"]
	page, valid := pageRaw.(int)
	if !valid {
		return nil, false, errors.New("invalid per_page data type")
	}

	perPageRaw := query["per_page"]
	perPage, valid = perPageRaw.(int)
	if !valid {
		return nil, false, errors.New("invalid per_page data type")
	}

	page = (page - 1) * perPage
	perPage++
	query["per_page"] = perPage
	query["page"] = page

	customers, err := cc.customer.SearchCustomer(ctx, query)
	if err != nil {
		return nil, false, err
	}

	// check for next page pagination
	hasNext := false
	if len(customers) >= perPage {
		hasNext = true
		customers = customers[:len(customers)-1]
	}

	return customers, hasNext, nil
}
