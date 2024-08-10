package pgsql

import (
	"context"
	"database/sql"
)

type pgsqlProvider interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type pgsqlRepository struct {
	pgsql pgsqlProvider
}

func NewPgsqlRepository(pgsql pgsqlProvider) *pgsqlRepository {
	return &pgsqlRepository{pgsql: pgsql}
}
