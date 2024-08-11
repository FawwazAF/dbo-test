package product

import (
	"context"

	"github.com/dbo-test/internal/model"
	"github.com/gin-gonic/gin"
)

type productControllerItf interface {
	GetAllProduct(ctx context.Context) ([]model.Product, error)
}

type responseWriter interface {
	GinHTTPResponseWriter(ctx *gin.Context, data interface{}, err error, httpStatus ...int)
}

type handler struct {
	product        productControllerItf
	responseWriter responseWriter
}

func NewHandler(product productControllerItf, writer responseWriter) *handler {
	return &handler{product: product, responseWriter: writer}
}
