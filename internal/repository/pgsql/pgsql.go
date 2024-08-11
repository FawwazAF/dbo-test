package pgsql

import (
	"context"
	"database/sql"
)

type pgsqlProvider interface {
	Begin() (*sql.Tx, error)
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

// SqlTx is the interface that wraps the SQL Tx methods.
type SqlTx interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func (repo *pgsqlRepository) BeginTx(ctx context.Context) (SqlTx, error) {
	return repo.pgsql.Begin()
}
