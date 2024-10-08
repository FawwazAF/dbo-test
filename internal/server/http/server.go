package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	jwt_pkg "github.com/dbo-test/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	Index    indexHandlerInterface
	Customer customerHandlerItf
	Order    orderHandlerItf
	Login    loginHandlerItf
	Product  productHandlerItf
}

type indexHandlerInterface interface {
	HandlerIndex(gCtx *gin.Context)
}

type customerHandlerItf interface {
	HandlerGetCustomerByID(g *gin.Context)
	HandlerAddCustomer(g *gin.Context)
	HandlerUpdateCustomer(g *gin.Context)
	HandlerDeleteCustomer(g *gin.Context)
	HandlerSearchCustomer(g *gin.Context)
}

type orderHandlerItf interface {
	HandlerGetOrderDetail(g *gin.Context)
	HandlerCreateOrder(g *gin.Context)
	HandlerDeleteOrder(g *gin.Context)
	HandlerUpdateOrder(g *gin.Context)
	HandlerSearchOrder(g *gin.Context)
}

type loginHandlerItf interface {
	HandlerLogin(c *gin.Context)
	HandlerLoginInfo(g *gin.Context)
}

type productHandlerItf interface {
	HandlerGetAllProduct(c *gin.Context)
}

type jwtMiddleware interface {
	GetJWTSecret() []byte
	ParseClientToken(tokenStr string) (*jwt.Token, error)
	GenerateJWTToken(param jwt_pkg.JWTTokenParameter) (string, error)
}

type Server struct {
	handler Handler
	router  *gin.Engine
	jwt     jwtMiddleware
}

func NewServer(handler Handler, jwt jwtMiddleware) *Server {
	return &Server{
		handler: handler,
		jwt:     jwt,
	}
}

func (s *Server) Start(addr string) {
	s.router = gin.Default()
	s.registerHandler()
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		log.Println("Running http server at port", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to run http server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
