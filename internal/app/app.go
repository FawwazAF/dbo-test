package app

import (
	// controller
	"fmt"
	"os"

	customer_ctrl "github.com/dbo-test/internal/controller/customer"
	login_ctrl "github.com/dbo-test/internal/controller/login"
	order_ctrl "github.com/dbo-test/internal/controller/order"
	product_ctrl "github.com/dbo-test/internal/controller/product"

	// repository
	"github.com/dbo-test/internal/repository/pgsql"

	// handler
	customer_handler "github.com/dbo-test/internal/server/http/customer"
	"github.com/dbo-test/internal/server/http/index"
	login_handler "github.com/dbo-test/internal/server/http/login"
	order_handler "github.com/dbo-test/internal/server/http/order"
	product_handler "github.com/dbo-test/internal/server/http/product"

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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbPGSQL := db_pgsql.NewDBSql("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName))
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
	productCtrl := product_ctrl.NewProduct(pgsqlRepo)

	// ===================================================== HANDLER ========================================================
	writer := http.NewHTTPWriter()
	handler := http.Handler{
		Index:    index.NewHandler(),
		Customer: customer_handler.NewHandler(customerCtrl, writer),
		Order:    order_handler.NewHandler(orderCtrl, writer),
		Login:    login_handler.NewHandler(loginCtrl, writer),
		Product:  product_handler.NewHandler(productCtrl, writer),
	}

	app.HTTPServers = http.NewServer(handler, jwt)
	return app, nil
}
