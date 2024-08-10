package app

import (
	// controller
	customer_ctrl "github.com/dbo-test/internal/controller/customer"

	// repository
	"github.com/dbo-test/internal/repository/pgsql"

	// handler
	customer_handler "github.com/dbo-test/internal/server/http/customer"
	"github.com/dbo-test/internal/server/http/index"

	"github.com/dbo-test/internal/server/http"
	"github.com/dbo-test/pkg/database/db_pgsql"
)

type Application struct {
	HTTPServers *http.Server
}

func NewApplication() (*Application, error) {
	// connect to db
	dbPGSQL := db_pgsql.NewDBSql("postgres", "user=admin password=admin1234 dbname=postgres host=localhost port=5432 sslmode=disable")
	pgsqlConn, err := dbPGSQL.ConnectSQLX()
	if err != nil {
		return nil, err
	}

	app := new(Application)

	// ===================================================== REPOSITORY =====================================================
	pgsqlRepo := pgsql.NewPgsqlRepository(pgsqlConn)

	// ===================================================== CONTROLLER =====================================================
	customerCtrl := customer_ctrl.NewCustomer(pgsqlRepo)

	// ===================================================== HANDLER ========================================================
	writer := http.NewHTTPWriter()
	handler := http.Handler{
		Index:    index.NewHandler(),
		Customer: customer_handler.NewHandler(customerCtrl, writer),
	}

	app.HTTPServers = http.NewServer(handler, nil)
	return app, nil
}
