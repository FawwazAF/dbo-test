package order

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dbo-test/internal/controller/order"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandlerGetOrderDetail(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	customerIDRaw, exist := g.Get("customer_id")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	customerID, valid := customerIDRaw.(float64)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer id from token"), http.StatusInternalServerError)
		return
	}

	orderIDRaw := g.Param("order_id")
	orderID, err := strconv.Atoi(orderIDRaw)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	order, err := h.order.GetOrderDetail(ctx, orderID, int(customerID))
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, order, nil, http.StatusOK)
}

func (h *handler) HandlerCreateOrder(g *gin.Context) {
	var (
		ctx = g.Request.Context()
		req OrderRequest
	)

	if err := g.ShouldBindJSON(&req); err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	orderReq := make([]order.OrderDetailRequest, len(req.OrderList))
	for i, v := range req.OrderList {
		orderReq[i] = order.OrderDetailRequest(v)
	}

	customerIDRaw, exist := g.Get("customer_id")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	customerID, valid := customerIDRaw.(float64)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer id from token"), http.StatusInternalServerError)
		return
	}

	err := h.order.CreateOrder(ctx, int(customerID), orderReq)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}

func (h *handler) HandlerDeleteOrder(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	customerIDRaw, exist := g.Get("customer_id")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	customerID, valid := customerIDRaw.(float64)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer id from token"), http.StatusInternalServerError)
		return
	}

	orderIDRaw := g.Param("order_id")
	orderID, err := strconv.Atoi(orderIDRaw)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	err = h.order.DeleteOrder(ctx, orderID, int(customerID))
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}

func (h *handler) HandlerUpdateOrder(g *gin.Context) {
	var (
		ctx = g.Request.Context()
		req UpdateOrderRequest
	)

	if err := g.ShouldBindJSON(&req); err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	orderReq := make([]order.OrderDetailRequest, len(req.OrderList))
	for i, v := range req.OrderList {
		orderReq[i] = order.OrderDetailRequest(v)
	}

	customerIDRaw, exist := g.Get("customer_id")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	customerID, valid := customerIDRaw.(float64)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer id from token"), http.StatusInternalServerError)
		return
	}

	err := h.order.UpdateOrder(ctx, req.OrderID, int(customerID), orderReq)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}

func (h *handler) HandlerSearchOrder(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	query, err := h.constructQuerySearchOrder(g)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	customerIDRaw, exist := g.Get("customer_id")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	customerID, valid := customerIDRaw.(float64)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer id from token"), http.StatusInternalServerError)
		return
	}

	metadata := map[string]interface{}{
		"filter_invoice": query["invoice"],
		"order_by":       query["order_by"],
		"page":           query["page"],
		"per_page":       query["per_page"],
	}

	orders, hasNext, err := h.order.SearchOrders(ctx, int(customerID), query)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}
	metadata["has_next"] = hasNext

	response := searchOrderResponse{
		OrderData: orders,
		Metadata:  metadata,
	}

	h.responseWriter.GinHTTPResponseWriter(g, response, nil, http.StatusOK)
}

func (h *handler) constructQuerySearchOrder(g *gin.Context) (map[string]interface{}, error) {
	query := make(map[string]interface{})
	invoice := g.Query("invoice")
	if invoice != "" {
		query["invoice"] = invoice
	}

	orderByRaw := g.Query("order_by")
	if orderByRaw != "" {
		orderBy := strings.ToLower(orderByRaw)
		if orderBy != "asc" && orderBy != "desc" {
			return nil, errors.New("invalid order_by")
		}

		query["order_by"] = orderBy
	}

	var (
		page    int
		perPage int
	)
	pageString := g.Query("page")
	if pageString != "" {
		page, _ = strconv.Atoi(pageString)
		query["page"] = page
	}

	perPageString := g.Query("per_page")
	if perPageString != "" {
		perPage, _ = strconv.Atoi(perPageString)
		query["per_page"] = perPage
	}

	return query, nil
}
