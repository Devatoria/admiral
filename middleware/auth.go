package middleware

import (
	"net/http"
	"strings"

	"github.com/Devatoria/admiral/auth"
	"github.com/Devatoria/admiral/models"

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

// AdminMiddleware is a Gin middleware ensuring HTTP basic auth authentication
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := auth.AuthenticateAdmin(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}

		c.Next()
	}
}

// ImageOwnerMiddleware checks that the user is the owner of the namespace containing the required image
func ImageOwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.GetCurrentUser(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		image := models.GetImageByName(strings.TrimPrefix(c.Param("image"), "/"))
		if image.ID == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if image.Namespace.Owner.ID != user.ID {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Keys["image"] = image
		c.Next()
	}
}

// ImageOwnerMiddleware add image in context
func ImageContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		image := models.GetImageByName(strings.TrimPrefix(c.Param("image"), "/"))
		if image.ID == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Keys["image"] = image
		c.Next()
	}
}
