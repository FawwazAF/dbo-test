package customer

import (
	"net/http"
	"strconv"

	"github.com/dbo-test/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandlerGetCustomerByID(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	idRaw := g.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
	}

	customer, err := h.customer.GetCustomerByID(ctx, id)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, customer, nil, http.StatusOK)
}

func (h *handler) HandlerAddCustomer(g *gin.Context) {
	var (
		ctx = g.Request.Context()
		req model.Customer
	)

	if err := g.BindJSON(&req); err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	err := h.customer.AddCustomer(ctx, &req)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}

func (h *handler) HandlerUpdateCustomer(g *gin.Context) {
	var (
		ctx = g.Request.Context()
		req model.Customer
	)

	if err := g.BindJSON(&req); err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	err := h.customer.UpdateCustomer(ctx, &req)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}

func (h *handler) HandlerDeleteCustomer(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	idRaw := g.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
	}

	err = h.customer.DeleteCustomer(ctx, id)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}
