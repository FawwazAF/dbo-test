package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware function to check JWT token and authorize user.
func (s *Server) authMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		r := c.Request
		w := c.Writer
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := s.jwt.ParseClientToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		exp, _ := claims["exp"].(time.Time)

		if hasTokenExpired(exp) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

func hasTokenExpired(exp time.Time) bool {
	return time.Now().After(exp)
}
