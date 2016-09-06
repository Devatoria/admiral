package api

import (
	"fmt"
	"net/http"

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

	r.GET("/version", getVersion)

	v1 := r.Group("/v1")
	{
		// Registry token endpoint
		v1.GET("/token", getToken)

		// Registry events notification endpoint
		v1.GET("/events", getEvents)
		v1.POST("/events", postEvents)

		// Namespace endpoints
		v1.GET("/namespaces", getNamespaces)
		v1.GET("/namespace/:id", getNamespace)
		v1.PUT("/namespace", putNamespace)
		v1.DELETE("/namespace/:id", deleteNamespace)
		v1.PATCH("/namespace/:id", patchNamespace)
		v1.GET("/namespace/:id/images", getNamespaceImages)

		// Image endpoints
		v1.GET("/images", getImages)
		v1.GET("/image/:id", getImage)
		v1.GET("/image/:id/tags", getImageTags)

		// User endpoints
		v1.GET("/users", getUsers)
		v1.GET("/user", getUser)
		v1.PUT("/user", putUser)
		v1.DELETE("/user/:id", deleteUser)
		v1.PATCH("/user/:id", patchUser)

		// Team endpoints
		v1.GET("/teams", getTeams)
		v1.GET("/team/:id", getTeam)
		v1.PUT("/team", putTeam)
		v1.DELETE("/team/:id", deleteTeam)
		v1.PATCH("/team/:id", patchTeam)

		// Team users management endpoints
		v1.GET("/teamUsers/:id", getTeamUsers)
		v1.POST("/teamUsers/:id", postTeamUsers)
		v1.DELETE("/teamUsers/:teamID/:userID", deleteTeamUser)
	}

	err := r.Run(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		panic(err)
	}
}
