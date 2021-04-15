package route

import (
	"github.com/catchup/registry-auth/auth"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	route := gin.New()

	authController := new(auth.AuthController)

	authRouteGroup := route.Group("/auth")
	{
		authRouteGroup.GET("/", authController.Test)
	}

	return route

}
