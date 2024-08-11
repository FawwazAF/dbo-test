package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandlerGetAllProduct(g *gin.Context) {
	var (
		ctx = g.Request.Context()
	)

	order, err := h.product.GetAllProduct(ctx)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, order, nil, http.StatusOK)
}
