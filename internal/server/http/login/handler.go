package login

import (
	"context"

	"github.com/gin-gonic/gin"
)

type loginControllerItf interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type responseWriter interface {
	GinHTTPResponseWriter(ctx *gin.Context, data interface{}, err error, httpStatus ...int)
}

type Handler struct {
	login          loginControllerItf
	responseWriter responseWriter
}

func NewHandler(login loginControllerItf, writer responseWriter) *Handler {
	return &Handler{login: login, responseWriter: writer}
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type response struct {
	Token string `json:"token"`
}

type LoginInfo struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}
