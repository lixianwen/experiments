package router

import (
	"gdemo/internal/handler/user"
	"gdemo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/signin", user.SignIn)
	r.POST("/signup", user.SignUp)
	r.POST("/refresh", middleware.Auth(), user.Refresh)

	// protected routes
	usergroup := r.Group("/api")
	usergroup.Use(middleware.Auth())
	{
		usergroup.GET("/users", user.GetAllUsers)
		usergroup.POST("/user", user.CreateUser)
		usergroup.GET("/user/:id", user.GetUser)
		usergroup.PUT("/user/:id", user.UpdateUser)
		usergroup.DELETE("/user/:id", user.DeleteUser)
	}

	return r
}
