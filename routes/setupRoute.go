package routes

import (
	"myapi/controller"

	"github.com/pilinux/gorest/config"
	"github.com/pilinux/gorest/lib/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(configure config.Configuration) (*gin.Engine, error) {
	if configure.Server.ServerEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.CORS(
		configure.Security.CORS.Origin,
		configure.Security.CORS.Credentials,
		configure.Security.CORS.Headers,
		configure.Security.CORS.Methods,
		configure.Security.CORS.MaxAge,
	))
	// For gorest <= v1.4.5
	// router.Use(middleware.CORS())

	// API:v1
	v1 := router.Group("/api/v1/")
	{
		// RDBMS
		if configure.Database.RDBMS.Activate == "yes" {
			// Register - no JWT required
			v1.POST("register", controller.CreateUserAuth)

			// Login - app issues JWT
			v1.POST("login", controller.Login)

			// Refresh - app issues new JWT
			rJWT := v1.Group("refresh")
			rJWT.Use(middleware.RefreshJWT())
			rJWT.POST("", controller.Refresh)

			// User
			rUsers := v1.Group("users")
			rUsers.GET("/:id", controller.GetUser) // Non-protected
			rUsers.Use(middleware.JWT())
			rUsers.POST("", controller.CreateUser) // Protected
			rUsers.PUT("", controller.UpdateUser)

			// Post
			rPosts := v1.Group("posts")
			rPosts.GET("", controller.GetPosts)    // Non-protected
			rPosts.GET("/:id", controller.GetPost) // Non-protected
			rPosts.Use(middleware.JWT())
			rPosts.POST("", controller.CreatePost) // Protected
		}
	}

	return router, nil
}
