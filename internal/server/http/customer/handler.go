package customer

import (
	"context"

	"github.com/dbo-test/internal/model"
	"github.com/gin-gonic/gin"
)

type customerControllerItf interface {
	GetCustomerByID(ctx context.Context, id int) (*model.Customer, error)
	AddCustomer(ctx context.Context, req *model.Customer) error
	UpdateCustomer(ctx context.Context, req *model.Customer) error
	DeleteCustomer(ctx context.Context, id int) error
}

type responseWriter interface {
	GinHTTPResponseWriter(ctx *gin.Context, data interface{}, err error, httpStatus ...int)
}

type handler struct {
	customer       customerControllerItf
	responseWriter responseWriter
}

func NewHandler(customer customerControllerItf, writer responseWriter) *handler {
	return &handler{customer: customer, responseWriter: writer}
}
