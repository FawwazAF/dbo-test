package customer

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dbo-test/internal/model"
	"github.com/gin-gonic/gin"
)

var (
	alphabetRegex = regexp.MustCompile(`[^a-zA-Z\s]+`)
	numericRegex  = regexp.MustCompile(`^\d+$`)
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
		return
	}

	err = h.customer.DeleteCustomer(ctx, id)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, "success", nil, http.StatusOK)
}

func (h *handler) HandlerSearchCustomer(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	query, err := h.constructQuerySearch(g)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}
	metadata := map[string]interface{}{
		"filter_name":         query["name"],
		"filter_phone_number": query["phone_number"],
		"order_by":            query["order_by"],
		"page":                query["page"],
		"per_page":            query["per_page"],
	}

	customers, hasNext, err := h.customer.SearchCustomer(ctx, query)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}
	metadata["has_next"] = hasNext

	response := searchCustomerResponse{
		CustomerData: customers,
		Metadata:     metadata,
	}

	h.responseWriter.GinHTTPResponseWriter(g, response, nil, http.StatusOK)
}

func (h *handler) constructQuerySearch(g *gin.Context) (map[string]interface{}, error) {
	query := make(map[string]interface{})
	name := g.Query("name")
	if name != "" {
		// remove invalid character for name
		cleanedName := alphabetRegex.ReplaceAllString(name, "")
		query["name"] = cleanedName
	}

	phoneNumber := g.Query("phone_number")
	if phoneNumber != "" {
		if !numericRegex.MatchString(phoneNumber) {
			return nil, errors.New("invalid phone number")
		}

		query["phone_number"] = phoneNumber
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
