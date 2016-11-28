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
			v1auth.GET("/image/*image", middleware.ImageOwnerMiddleware(), getImage)
			v1auth.DELETE("/image/*image", middleware.ImageOwnerMiddleware(), deleteImage)
			v1auth.PATCH("/image/public/*image", middleware.ImageOwnerMiddleware(), setImagePublic)
			v1auth.PATCH("/image/private/*image", middleware.ImageOwnerMiddleware(), setImagePrivate)

			v1auth.GET("/login", getLogin)
			v1auth.GET("/token", getToken)
		}
		
		v1admin := r.Group("/v1/admin")
		v1admin.Use(middleware.AdminMiddleware())
		{
			v1admin.GET("/images", getAllImages)
			v1admin.DELETE("/image/*image", middleware.ImageContextMiddleware(), deleteImage)
			v1admin.GET("/image/*image", middleware.ImageContextMiddleware(), getImage)
			v1admin.PATCH("/image/public/*image", middleware.ImageContextMiddleware(), setImagePublic)
			v1admin.PATCH("/image/private/*image", middleware.ImageContextMiddleware(), setImagePrivate)

			v1admin.GET("/login", getLogin)
		}
	}

	err := r.Run(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		panic(err)
	}
}
