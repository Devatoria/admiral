package middleware

import (
	"net/http"

	"github.com/Devatoria/admiral/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware ensuring HTTP basic auth authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.Authenticate(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}

		c.Keys["user"] = user
		c.Next()
	}
}
