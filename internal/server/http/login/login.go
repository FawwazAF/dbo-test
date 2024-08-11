package login

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleLogin(g *gin.Context) {
	var (
		req LoginRequest
		ctx = g.Request.Context()
	)

	if err := g.ShouldBindJSON(&req); err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusBadRequest)
		return
	}

	token, err := h.login.Login(ctx, req.Username, req.Password)
	if err != nil {
		h.responseWriter.GinHTTPResponseWriter(g, nil, err, http.StatusInternalServerError)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, response{
		Token: token,
	}, nil, http.StatusOK)
}

func (h *Handler) HandleLoginInfo(g *gin.Context) {
	idRaw, exist := g.Get("customer_id")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("you are not logged in"), http.StatusUnauthorized)
		return
	}

	id, valid := idRaw.(float64)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer id from token"), http.StatusUnauthorized)
		return
	}

	usernameRaw, exist := g.Get("username")
	if !exist {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("you are not logged in"), http.StatusUnauthorized)
		return
	}

	username, valid := usernameRaw.(string)
	if !valid {
		h.responseWriter.GinHTTPResponseWriter(g, nil, errors.New("got invalid customer username from token"), http.StatusUnauthorized)
		return
	}

	h.responseWriter.GinHTTPResponseWriter(g, LoginInfo{ID: int(id), Username: username}, nil, http.StatusOK)
}
