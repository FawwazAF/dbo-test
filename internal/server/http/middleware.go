package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Option int

const (
	SuperAdminAccess Option = iota + 1
)

// Middleware function to check JWT token and authorize user.
func (s *Server) authMiddleware(opts ...Option) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := s.jwt.ParseClientToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		usernameRaw := claims["username"]
		for _, opt := range opts {
			if opt == SuperAdminAccess {
				username, valid := usernameRaw.(string)
				if !valid {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "got invalid username"})
					c.Abort()
					return
				}
				if username != "superadmin" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Need Super Admin Access!!"})
					c.Abort()
					return
				}
			}
		}

		c.Set("customer_id", claims["id"])
		c.Set("username", usernameRaw)
		c.Next()
	})
}
