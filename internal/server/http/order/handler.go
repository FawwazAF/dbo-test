package order

import (
	"context"

	"github.com/dbo-test/internal/controller/order"
	"github.com/dbo-test/internal/model"
	"github.com/gin-gonic/gin"
)

type orderControllerItf interface {
	GetOrderDetail(ctx context.Context, orderID, customerID int) (*model.Order, error)
	CreateOrder(ctx context.Context, customerID int, req []order.OrderDetailRequest) error
	DeleteOrder(ctx context.Context, orderID, customerID int) error
	UpdateOrder(ctx context.Context, orderID, customerID int, req []order.OrderDetailRequest) error
	SearchOrders(ctx context.Context, customerID int, queryParameters map[string]interface{}) ([]model.Order, bool, error)
}

type responseWriter interface {
	GinHTTPResponseWriter(ctx *gin.Context, data interface{}, err error, httpStatus ...int)
}

type handler struct {
	order          orderControllerItf
	responseWriter responseWriter
}

func NewHandler(order orderControllerItf, writer responseWriter) *handler {
	return &handler{order: order, responseWriter: writer}
}

type OrderRequest struct {
	OrderList []OrderDetailRequest `json:"order_list"`
}

type OrderDetailRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type UpdateOrderRequest struct {
	OrderID   int                  `json:"order_id"`
	OrderList []OrderDetailRequest `json:"order_list"`
}

type searchOrderResponse struct {
	OrderData []model.Order          `json:"order_data"`
	Metadata  map[string]interface{} `json:"metadata"`
}
