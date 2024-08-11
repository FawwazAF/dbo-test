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

func (repo *pgsqlRepository) GetCustomerByID(ctx context.Context, id int) (*model.Customer, error) {
	query := `
		SELECT 
			id,
			username,
			name,
			email,
			phone_number,
			date_of_birth,
			address,
			status,
			created_at,
			updated_at
		FROM
			dbo_trx_customer
		WHERE
			id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	customer := model.Customer{}
	if err := repo.pgsql.GetContext(ctx, &customer, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	return &customer, nil
}

func (repo *pgsqlRepository) GetCustomerByUsername(ctx context.Context, username string) (*model.Customer, error) {
	query := `
		SELECT 
			id,
			password
		FROM
			dbo_trx_customer
		WHERE
			username = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	customer := model.Customer{}
	if err := repo.pgsql.GetContext(ctx, &customer, query, username); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &customer, nil
}

func (repo *pgsqlRepository) AddCustomer(ctx context.Context, req *model.Customer) error {
	query := `
		INSERT INTO dbo_trx_customer(
			username,
			password,
			status,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5)
	`
	args := []interface{}{
		req.Username,
		string(req.GetHashedPassword()),
		req.Status,
		time.Now(),
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	_, err := repo.pgsql.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) UpdateCustomer(ctx context.Context, req *model.Customer) error {
	query, args := repo.buildUpdateCustomerQuery(req)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	_, err := repo.pgsql.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) buildUpdateCustomerQuery(req *model.Customer) (string, []interface{}) {
	query := `UPDATE dbo_trx_customer SET `
	builder := strings.Builder{}
	args := []interface{}{}
	builder.WriteString(query)

	if req.Name != "" {
		builder.WriteString(`name=?, `)
		args = append(args, req.Name)
	}

	if req.Email != "" {
		builder.WriteString(`email=?, `)
		args = append(args, req.Email)
	}

	if req.PhoneNumber != "" {
		builder.WriteString(`phone_number=?, `)
		args = append(args, req.PhoneNumber)
	}

	if !req.DateOfBirth.IsZero() {
		builder.WriteString(`date_of_birth=?, `)
		args = append(args, req.DateOfBirth)
	}

	if req.Address != "" {
		builder.WriteString(`address=?, `)
		args = append(args, req.Address)
	}

	builder.WriteString(`updated_at=? WHERE id = ?`)
	args = append(args, time.Now(), req.ID)

	return sqlx.Rebind(sqlx.DOLLAR, builder.String()), args
}

func (repo *pgsqlRepository) DeleteCustomer(ctx context.Context, id int) error {
	query := `
		DELETE FROM dbo_trx_customer
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	_, err := repo.pgsql.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *pgsqlRepository) SearchCustomer(ctx context.Context, parameters map[string]interface{}) ([]model.Customer, error) {
	query, args := repo.buildSearchCustomerQuery(parameters)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	customers := []model.Customer{}
	if err := repo.pgsql.SelectContext(ctx, &customers, query, args...); err != nil {
		return nil, err
	}

	return customers, nil
}

func (repo *pgsqlRepository) buildSearchCustomerQuery(parameters map[string]interface{}) (string, []interface{}) {
	query := `
		SELECT 
			id,
			username,
			name,
			email,
			phone_number,
			date_of_birth,
			address,
			status,
			created_at,
			updated_at
		FROM
			dbo_trx_customer 
	`
	builder := strings.Builder{}
	args := []interface{}{}
	builder.WriteString(query)

	whereClause := []string{}
	if nameValue := parameters["name"]; nameValue != nil {
		whereClause = append(whereClause, `"name" LIKE (?) `)
		args = append(args, "%"+fmt.Sprintf("%v", nameValue)+"%")
	}

	if phoneNumber := parameters["phone_number"]; phoneNumber != nil {
		whereClause = append(whereClause, "phone_number LIKE (?) ")
		args = append(args, fmt.Sprintf("%v", phoneNumber)+"%")
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
