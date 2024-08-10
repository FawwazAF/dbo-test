package db_pgsql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBSQL struct {
	DriverName string
	Source     string
}

func NewDBSql(driver, source string) *DBSQL {
	return &DBSQL{
		DriverName: driver,
		Source:     source,
	}
}

func (d *DBSQL) ConnectSQLX() (*sqlx.DB, error) {
	db, err := sqlx.Connect(d.DriverName, d.Source)
	if err != nil {
		return nil, err
	}

	return db, nil
}
