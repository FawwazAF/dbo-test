package app

import (
	// controller
	customer_ctrl "github.com/dbo-test/internal/controller/customer"
	login_ctrl "github.com/dbo-test/internal/controller/login"
	order_ctrl "github.com/dbo-test/internal/controller/order"

	// repository
	"github.com/dbo-test/internal/repository/pgsql"

	// handler
	customer_handler "github.com/dbo-test/internal/server/http/customer"
	"github.com/dbo-test/internal/server/http/index"
	login_handler "github.com/dbo-test/internal/server/http/login"
	order_handler "github.com/dbo-test/internal/server/http/order"

	"github.com/dbo-test/internal/server/http"
	"github.com/dbo-test/pkg/database/db_pgsql"
	"github.com/dbo-test/pkg/jwt"
)

// Secret key for signing tokens
var secretKey = []byte("testing1234567890!@#$%^&*()")

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
	// init jwt
	jwt := jwt.NewJWT(string(secretKey))
	// ===================================================== REPOSITORY =====================================================
	pgsqlRepo := pgsql.NewPgsqlRepository(pgsqlConn)

	// ===================================================== CONTROLLER =====================================================
	customerCtrl := customer_ctrl.NewCustomer(pgsqlRepo)
	orderCtrl := order_ctrl.NewOrder(pgsqlRepo)
	loginCtrl := login_ctrl.NewLogin(pgsqlRepo, jwt)

	// ===================================================== HANDLER ========================================================
	writer := http.NewHTTPWriter()
	handler := http.Handler{
		Index:    index.NewHandler(),
		Customer: customer_handler.NewHandler(customerCtrl, writer),
		Order:    order_handler.NewHandler(orderCtrl, writer),
		Login:    login_handler.NewHandler(loginCtrl, writer),
	}

	app.HTTPServers = http.NewServer(handler, jwt)
	return app, nil
}
