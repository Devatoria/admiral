package middleware

import (
	"net/http"

	"github.com/Devatoria/admiral/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware ensuring HTTP basic auth authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		if err = auth.Authenticate(c.Request); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
