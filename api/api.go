package api

import (
	"fmt"
	"net/http"

	"github.com/Devatoria/admiral/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Run func registers endpoints and runs the API
func Run(address string, port int) {
	if !viper.GetBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/events", postEvents)

	v1 := r.Group("/v1")
	{
		v1.GET("/version", getVersion)
		v1.PUT("/user", putUser)

		// Authenticated endpoints
		v1auth := r.Group("/v1")
		v1auth.Use(middleware.AuthMiddleware())
		{
			v1auth.GET("/images", getImages)
			v1auth.DELETE("/image/*image", middleware.ImageOwnerMiddleware(), deleteImage)
			v1auth.PATCH("/image/public/*image", middleware.ImageOwnerMiddleware(), setImagePublic)
			v1auth.PATCH("/image/private/*image", middleware.ImageOwnerMiddleware(), setImagePrivate)

			v1auth.GET("/login", getLogin)
			v1auth.GET("/token", getToken)
		}
	}

	err := r.Run(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		panic(err)
	}
}
