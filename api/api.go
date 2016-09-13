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
			// Registry token endpoint
			v1auth.GET("/token", getToken)

			// Registry events notifications endpoint
			v1auth.GET("/events", getEvents)

			// Namespace endpoints
			v1auth.GET("/namespaces", getNamespaces)
			v1auth.GET("/namespace/:id", getNamespace)
			v1auth.PUT("/namespace", putNamespace)
			v1auth.DELETE("/namespace/:id", deleteNamespace)
			v1auth.PATCH("/namespace/:id", patchNamespace)
			v1auth.GET("/namespace/:id/images", getNamespaceImages)

			// Image endpoints
			v1auth.GET("/images", getImages)
			v1auth.GET("/image/:id", getImage)
			v1auth.GET("/image/:id/tags", getImageTags)

			// User endpoints
			v1auth.GET("/users", getUsers)
			v1auth.GET("/user", getUser)
			v1auth.DELETE("/user/:id", deleteUser)
			v1auth.PATCH("/user/:id", patchUser)

			// Team endpoints
			v1auth.GET("/teams", getTeams)
			v1auth.GET("/team/:id", getTeam)
			v1auth.PUT("/team", putTeam)
			v1auth.DELETE("/team/:id", deleteTeam)
			v1auth.PATCH("/team/:id", patchTeam)

			// Team users management endpoints
			v1auth.GET("/teamUsers/:id", getTeamUsers)
			v1auth.POST("/teamUsers/:id", postTeamUsers)
			v1auth.DELETE("/teamUsers/:teamID/:userID", deleteTeamUser)

			// Team namespace rights management
			v1auth.GET("/teamNamespaceRights/:teamID/:namespaceID", getTeamNamespaceRights)
			v1auth.PUT("/teamNamespaceRight", putTeamNamespaceRight)
			v1auth.DELETE("/teamNamespaceRight/:teamID/:namespaceID", deleteTeamNamespaceRight)
		}
	}

	err := r.Run(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		panic(err)
	}
}
